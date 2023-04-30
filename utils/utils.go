package utils

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetChatId(update tgbotapi.Update) int64 {
	return update.Message.Chat.ID
}
