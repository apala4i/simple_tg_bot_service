package tasks

import (
	"testProjects/priceChecker/simple_tg_bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TaskImpl struct {
	ActionFunc func(tgBot *services.TgBot, update tgbotapi.Update) error
}

func (t *TaskImpl) Action(tgBot *services.TgBot, update tgbotapi.Update) error {
	return t.ActionFunc(tgBot, update)
}

func NewTask(ActionFunc func(tgBot *services.TgBot, update tgbotapi.Update) error) services.Task {
	return &TaskImpl{ActionFunc: ActionFunc}
}
