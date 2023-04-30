package services

import (
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

const (
	baseTasksSize = 10
	prefix        = "@"
)

type TgBotServer struct {
	bot      *TgBot
	tasks    []Task
	infOnErr Task
}

type Task interface {
	Action(tgBot *TgBot, update tgbotapi.Update) error
	GetNamePattern() *regexp.Regexp
	CompareName(string) bool
	GetDescription() string
}

type TaskWithData interface {
	Action(tgBot *TgBot, chatId int64, data struct{}) error
}

// noinspection GoUnusedExportedFunction
func NewTgBotServer(bot *TgBot) *TgBotServer {
	return &TgBotServer{bot: bot, tasks: make([]Task, 0, baseTasksSize)}
}

// noinspection GoUnusedExportedFunction
func NewBaseTgBotServer(token string) *TgBotServer {
	tgBotApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	tgBotApi.Debug = true
	bot := NewTgBot(tgBotApi)
	return NewTgBotServer(bot)
}

func (c *TgBotServer) AddTask(newTask Task) bool {
	index := slices.IndexFunc(c.tasks, func(task Task) bool {
		return task.CompareName(newTask.GetNamePattern().String())
	})
	if index == -1 {
		c.tasks = append(c.tasks, newTask)
		return true
	}
	logrus.Infof("[AddTask] task with such name pattern alreadyExists. taskName: %v", newTask.GetNamePattern().String())
	return false
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
		command, ok := c.isValidCommand(update.Message.Text)
		if !ok {
			continue
		}
		index := slices.IndexFunc(c.tasks, func(task Task) bool {
			return task.CompareName(command)
		})
		if index != -1 {
			err := c.tasks[index].Action(c.bot, update)
			if err != nil {
				logrus.Errorf("error, while processing task: %v", err)
			}
		} else if c.infOnErr != nil {
			err := c.infOnErr.Action(c.bot, update)
			if err != nil {
				logrus.Errorf("error, while processing task: %v", err)
			}
		}

	}
}

func (c *TgBotServer) isValidCommand(command string) (string, bool) {
	sl := strings.Split(command, " ")
	if len(sl) > 2 && sl[len(sl)-1] == prefix+c.bot.GetBot().Self.UserName {
		return strings.Join(sl[:len(sl)-1], " "), true
	}
	return "", false
}

func (c *TgBotServer) GetTasks() []Task {
	return c.tasks
}

func (c *TgBotServer) EnableInfo(task Task) {
	c.infOnErr = task
}
