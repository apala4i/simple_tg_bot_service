package utils

import (
	"strings"

	"github.com/apala4i/simple_tg_bot_service/services"
	"github.com/apala4i/simple_tg_bot_service/tasks"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetInfoTask(server *services.TgBotServer) services.Task {
	return tasks.NewTask(func(tgBot *services.TgBot, update tgbotapi.Update) error {
		resStr := strings.Builder{}
		tasksSl := server.GetTasks()
		for i := range tasksSl {
			resStr.WriteString(tasksSl[i].GetDescription() + "\n")
		}
		return tgBot.SendMessage(GetChatId(update), resStr.String())
	}, "/info", "show available commands")
}

func GetCustomInfoTask(msg string, dsc string) services.Task {
	return tasks.NewTask(func(tgBot *services.TgBot, update tgbotapi.Update) error {
		return tgBot.SendMessage(GetChatId(update), msg)
	}, "/info", dsc)
}
