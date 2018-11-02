package telegram

import (
	"net/http"
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
