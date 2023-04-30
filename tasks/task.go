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

type taskImpl struct {
	actionFunc  func(tgBot *services.TgBot, update tgbotapi.Update) error
	namePattern *regexp.Regexp
	description string
}

func (t *taskImpl) Action(tgBot *services.TgBot, update tgbotapi.Update) error {
	return t.actionFunc(tgBot, update)
}

func (t *taskImpl) GetNamePattern() *regexp.Regexp {
	return t.namePattern
}

func (t *taskImpl) CompareName(name string) bool {
	return t.GetNamePattern().Match([]byte(name))
}

// noinspection GoUnusedExportedFunction
func NewTask(actionFunc func(tgBot *services.TgBot,
	update tgbotapi.Update) error,
	name string,
	dsc string) services.Task {
	task := &taskImpl{actionFunc: actionFunc, description: dsc}
	pattern, err := regexp.Compile(name)
	if err != nil {
		logrus.Errorf("regexp compile error: %v", err)
	}
	task.namePattern = pattern
	return task
}

func (t *taskImpl) GetDescription() string {
	if len(t.description) > 0 {
		return t.GetNamePattern().String() + " - " + t.description
	}
	return ""
}
