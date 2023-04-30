package main

import (
	"time"

	"github.com/apala4i/simple_tg_bot_service/maintainer"
	"github.com/apala4i/simple_tg_bot_service/services"
	"github.com/apala4i/simple_tg_bot_service/tasks"
	"github.com/apala4i/simple_tg_bot_service/utils"
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

	cronTask := tasks.NewDefaultCronTask(tasks.NewTask(func(tgBot *services.TgBot, update tgbotapi.Update) error {
		return tgBot.SendMessage(utils.GetChatId(update), "cronTask")
	}, "/cronTask"), time.Second*5)

	// add tasks to tg bot
	bot.AddTask(task)
	bot.AddTask(cronTask)

	mt.AddServer("bot name", bot)

	mt.Start()
}
