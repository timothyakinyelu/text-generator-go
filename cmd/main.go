package main

import (
	"log"
	new_bot "text-generator/src"
	"time"

	"github.com/NicoNex/echotron/v3"
)

func main() {
	echotron.NewAPI(new_bot.TelegramToken).SetMyCommands(nil, new_bot.Commands...)

	dsp := echotron.NewDispatcher(new_bot.TelegramToken, new_bot.NewBot)
	for {
		err := dsp.Poll()
		if err != nil {
			log.Println("Error polling updates:", err)
		}
		time.Sleep(5 * time.Second)
	}
}
