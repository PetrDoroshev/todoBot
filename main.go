package main

import (
	"github.com/PetrDoroshev/todoBot/bot"
	"github.com/PetrDoroshev/todoBot/data_base"
)

func main() {

	db := data_base.NewDB()
	b := bot.NewBot(db)
	b.BotLoop()
	db.CloseDB()
}
