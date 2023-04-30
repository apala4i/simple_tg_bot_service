package main

import (
	"fmt"
	"strconv"
	"strings"
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
	bot := services.NewBaseTgBotServer("test")

	// create task
	task := tasks.NewTask(func(tgBot *services.TgBot, update tgbotapi.Update) error {
		digit, _ := strconv.Atoi(strings.Split(utils.GetMsgText(update), " ")[1])

		return tgBot.SendMessage(update.Message.Chat.ID, fmt.Sprintf("got %v", digit))
	}, "/hello \\d+", "send hello")

	cronTask := tasks.NewDefaultCronTask(tasks.NewTask(func(tgBot *services.TgBot, update tgbotapi.Update) error {
		return tgBot.SendMessage(utils.GetChatId(update), "cronTask")
	}, "/cronTask", "send cronTask every 5 seconds"), time.Second*5)

	// add tasks to tg bot
	bot.AddTask(task)
	bot.AddTask(cronTask)
	bot.EnableInfo(utils.GetInfoTask(bot))

	mt.AddServer("bot name", bot)

	mt.Start()
}
