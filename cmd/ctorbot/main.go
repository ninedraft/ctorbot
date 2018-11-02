package main

import (
	"context"
	"log"

	"github.com/ninedraft/ctorbot/pkg/telegram"
)

func main() {
	var bot, errNewBot = telegram.New(telegram.BotConfig{})
	if errNewBot != nil {
		log.Fatal(errNewBot)
	}
	if err := bot.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
