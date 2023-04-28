package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type TgBotServer struct {
	bot   *TgBot
	tasks map[string]Task
}

type Task interface {
	Action(tgBot *TgBot, update tgbotapi.Update) error
}

type TaskWithData interface {
	Action(tgBot *TgBot, chatId int64, data struct{}) error
}

// noinspection GoUnusedExportedFunction
func NewTgBotServer(bot *TgBot) *TgBotServer {
	return &TgBotServer{bot: bot, tasks: make(map[string]Task)}
}

func (c *TgBotServer) AddTask(endpoint string, task Task) {
	c.tasks[endpoint] = task
}

func (c *TgBotServer) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.bot.Api.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}
		task, ok := c.tasks[update.Message.Text]
		if ok {
			err := task.Action(c.bot, update)
			if err != nil {
				logrus.Errorf("error, while processing task: %v", err)
			}
		}

	}
}