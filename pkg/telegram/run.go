package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/hekmon/transmissionrpc"
	"gopkg.in/telegram-bot-api.v4"
)

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
