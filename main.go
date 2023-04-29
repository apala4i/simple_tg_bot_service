package main

import (
	"github.com/apala4i/simple_tg_bot_service/maintainer"
	"github.com/apala4i/simple_tg_bot_service/services"
	"github.com/apala4i/simple_tg_bot_service/tasks"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// simple example of library usage
func main() {

	// create maintainer for tg bots
	mt := maintainer.NewMaintainer()

	// create tg bot
	bot := services.NewBaseTgBotServer("put_your_token_here")

	// create task
	task := tasks.NewTask(func(tgBot *services.TgBot, update tgbotapi.Update) error {
		return tgBot.SendMessage(update.Message.Chat.ID, "hello")
	}, "/hello")

	// add task to tg bot
	bot.AddTask(task)

	mt.AddServer("bot name", bot)

	mt.Start()
}
