package main

import (
	"github.com/apala4i/simple_tg_bot_service/maintainer"
	"github.com/apala4i/simple_tg_bot_service/services"
	"github.com/apala4i/simple_tg_bot_service/tasks"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	maintainer := maintainer.NewMaintainer()
	bot, err := tgbotapi.NewBotAPI("testtoken")
	if err != nil {
		logrus.Panicf("cannot init Bot. Error: %s", err)
	}
	server := services.NewTgBotServer(services.NewTgBot(bot))
	server.AddTask("/popa", tasks.NewTask(
		func(tgBot *services.TgBot, update tgbotapi.Update) error {
			err := tgBot.SendMessage(update.Message.Chat.ID, "popa")
			if err != nil {
				return err
			}
			return nil
		}))

	maintainer.AddServer(bot.Self.UserName, server)
	maintainer.Start()
}
