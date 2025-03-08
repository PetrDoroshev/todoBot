package main

import (
	"github.com/PetrDoroshev/HomeWork_db/golang/todo_list/bot"
	_ "github.com/PetrDoroshev/HomeWork_db/golang/todo_list/bot"
	"github.com/PetrDoroshev/HomeWork_db/golang/todo_list/data_base"
)

func main() {

	db := data_base.NewDB()
	b := bot.NewBot(db)
	b.BotLoop()
	db.CloseDB()
}
