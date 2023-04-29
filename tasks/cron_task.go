package tasks

import (
	"sync"
	"time"

	"github.com/apala4i/simple_tg_bot_service/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type CronTask struct {
	running     map[int64]struct{}
	m           sync.Mutex
	task        services.Task
	failureTask services.Task
	stopTask    services.Task
	sleepTime   time.Duration
}

type BaseFailureTask struct {
	msg string
}

const defaultMsg = "already running"

func newBaseFailureTask(msg string) services.Task {
	return NewTask(func(tgBot *services.TgBot, update tgbotapi.Update) error {
		return tgBot.ReplyMessage(update, msg)
	}, baseName)
}

func NewCronTask(task services.Task, failureTask services.Task, sleepTime time.Duration) services.Task {
	return &CronTask{running: make(map[int64]struct{}), task: task, failureTask: failureTask, sleepTime: sleepTime}
}

func NewReplyCronAction(task services.Task, msg string, sleepTime time.Duration) services.Task {
	return NewCronTask(task, newBaseFailureTask(msg), sleepTime)
}

// noinspection GoUnusedExportedFunction
func NewDefaultCronAction(task services.Task, sleepTime time.Duration) services.Task {
	return NewReplyCronAction(task, defaultMsg, sleepTime)
}

func (c *CronTask) Action(tgBot *services.TgBot, update tgbotapi.Update) error {
	_, ok := c.running[update.Message.Chat.ID]
	if ok && c.failureTask.CompareName(update.Message.Text) {
		c.deleteFromQueue(update.Message.Chat.ID)
		err := c.failureTask.Action(tgBot, update)
		if err != nil {
			logrus.Errorf("stop action failed. error: %v", err)
		}
	}
	if !ok {
		c.running[update.Message.Chat.ID] = struct{}{}
		go func(running *map[int64]struct{}) {
			for {
				_, ok := (*running)[update.Message.Chat.ID]
				if !ok {
					return
				}
				err := c.task.Action(tgBot, update)
				if err != nil {
					logrus.Errorf("[CronAction.Action] error: %v", err)
					c.deleteFromQueue(update.Message.Chat.ID)
					return
				}
				time.Sleep(c.sleepTime)
			}
		}(&c.running)
	} else {
		return c.failureTask.Action(tgBot, update)
	}
	return nil
}

func (c *CronTask) deleteFromQueue(chatId int64) {
	c.m.Lock()
	delete(c.running, chatId)
	c.m.Unlock()
}

func (c *CronTask) GetName() string {
	return c.task.GetName()
}

func (c *CronTask) CompareName(name string) bool {
	return c.task.GetName() == name || c.stopTask.GetName() == name
}

func (c *CronTask) addPostfixToName(s string) {

}
