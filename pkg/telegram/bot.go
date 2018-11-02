package telegram

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/hekmon/transmissionrpc"
	"golang.org/x/net/proxy"
	"gopkg.in/telegram-bot-api.v4"
)

type BotConfig struct {
	TelegramToken string
	ProxyConfig
}

type ProxyConfig struct {
	ProxyAddress  string
	ProxyUsername string
	ProxyPassword string
}

func (pconfig ProxyConfig) IsEmpty() bool {
	return pconfig == (ProxyConfig{})
}

func (pconfig ProxyConfig) ProxyAuth() *proxy.Auth {
	return &proxy.Auth{
		User:     pconfig.ProxyUsername,
		Password: pconfig.ProxyPassword,
	}
}

type Bot struct {
	transmission *transmissionrpc.Client
	tg           *tgbotapi.BotAPI
}

func New(config BotConfig) (Bot, error) {
	var torrentClient, err = transmissionrpc.New(
		"localhost", // host
		"",          // username
		"",          // password
		&transmissionrpc.AdvancedConfig{
			HTTPTimeout: 10 * time.Second,
		})
	if err != nil {
		return Bot{}, err
	}
	var tgClient, newBotErr = tgbotapi.NewBotAPI(config.TelegramToken)
	if newBotErr != nil {
		return Bot{}, newBotErr
	}
	if !config.IsEmpty() {
		var dialer, errSOCKS5 = proxy.SOCKS5("tcp",
			config.ProxyAddress,
			config.ProxyAuth(),
			proxy.Direct)
		if errSOCKS5 != nil {
			return Bot{}, errSOCKS5
		}
		tgClient.Client.Transport = &http.Transport{
			Dial: dialer.Dial,
		}
	}
	return Bot{
		transmission: torrentClient,
		tg:           tgClient,
	}, nil
}

func (bot Bot) Run(ctx context.Context) error {
	var errCh = make(chan error)
	var end, stop = stopper()
	go func() {
		defer stop()
		var config = tgbotapi.NewUpdate(0)
		var updates, getUpdErr = bot.tg.GetUpdatesChan(config)
		if getUpdErr != nil {
			errCh <- getUpdErr
			return
		}
		for {
			select {
			case <-ctx.Done():
				return
			case <-end:
				return
			case upd := <-updates:
				switch {
				case upd.Message != nil:
					var msg = upd.Message
					log.Printf("%s: %q\n", msg.From.UserName, msg.Text)
					if msg.IsCommand() {
						switch msg.Command() {
						case "list":
						case "add":
							var magnet = strings.TrimSpace(msg.CommandArguments())
							var torrent, err = bot.transmission.TorrentAdd(&transmissionrpc.TorrentAddPayload{
								Filename: &magnet,
							})
							if err != nil {
								log.Println(err)
								bot.tg.Send(tgbotapi.NewMessage(msg.Chat.ID, err.Error()))
								continue
							}
							bot.tg.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("torrent %q is added", *torrent.Name)))
						}
					}
				default:
					continue
				}
			}
		}
	}()
	select {
	case err := <-errCh:
		log.Fatal(err)
	case <-ctx.Done():
		return nil
	}
	return nil
}

func stopper() (<-chan struct{}, func()) {
	var stopCh = make(chan struct{})
	var oncer = &sync.Once{}
	var stopFn = func() {
		oncer.Do(func() {
			close(stopCh)
		})
	}
	return stopCh, stopFn
}
