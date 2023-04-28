package tasks

import (
	"github.com/apala4i/simple_tg_bot_service/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TaskImpl struct {
	ActionFunc func(tgBot *services.TgBot, update tgbotapi.Update) error
}

func (t *TaskImpl) Action(tgBot *services.TgBot, update tgbotapi.Update) error {
	return t.ActionFunc(tgBot, update)
}

// noinspection GoUnusedExportedFunction
func NewTask(ActionFunc func(tgBot *services.TgBot, update tgbotapi.Update) error) services.Task {
	return &TaskImpl{ActionFunc: ActionFunc}
}
