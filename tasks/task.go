package tasks

import (
	"github.com/apala4i/simple_tg_bot_service/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const baseName = "baseName"

type TaskImpl struct {
	ActionFunc func(tgBot *services.TgBot, update tgbotapi.Update) error
	name       string
}

func (t *TaskImpl) Action(tgBot *services.TgBot, update tgbotapi.Update) error {
	return t.ActionFunc(tgBot, update)
}

func (t *TaskImpl) GetName() string {
	return t.name
}

func (t *TaskImpl) CompareName(name string) bool {
	return t.GetName() == name
}

// noinspection GoUnusedExportedFunction
func NewTask(actionFunc func(tgBot *services.TgBot, update tgbotapi.Update) error, name string) services.Task {
	return &TaskImpl{ActionFunc: actionFunc, name: name}
}

func NewUnnamedTask(actionFunc func(tgBot *services.TgBot, update tgbotapi.Update) error) services.Task {
	return NewTask(actionFunc, baseName)
}
