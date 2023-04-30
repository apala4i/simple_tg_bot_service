package tasks

import (
	"regexp"

	"github.com/apala4i/simple_tg_bot_service/services"
	"github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	baseName = "baseName"
)

type TaskImpl struct {
	ActionFunc  func(tgBot *services.TgBot, update tgbotapi.Update) error
	namePattern *regexp.Regexp
}

func (t *TaskImpl) Action(tgBot *services.TgBot, update tgbotapi.Update) error {
	return t.ActionFunc(tgBot, update)
}

func (t *TaskImpl) GetNamePattern() *regexp.Regexp {
	return t.namePattern
}

func (t *TaskImpl) CompareName(name string) bool {
	return t.GetNamePattern().Match([]byte(name))
}

// noinspection GoUnusedExportedFunction
func NewTask(actionFunc func(tgBot *services.TgBot, update tgbotapi.Update) error, name string) services.Task {
	task := &TaskImpl{ActionFunc: actionFunc}
	pattern, err := regexp.Compile(name)
	if err != nil {
		logrus.Errorf("regexp compile error: %v", err)
	}
	task.namePattern = pattern
	return task

}

// noinspection GoUnusedExportedFunction
func NewUnnamedTask(actionFunc func(tgBot *services.TgBot, update tgbotapi.Update) error) services.Task {
	return NewTask(actionFunc, baseName)
}
