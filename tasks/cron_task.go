package tasks

import (
	"sync"
	"testProjects/priceChecker/simple_tg_bot/services"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type CronTask struct {
	running     map[int64]struct{}
	m           sync.Mutex
	task        services.Task
	failureTask services.Task
	sleepTime   time.Duration
}

type BaseFailureTask struct {
	msg string
}

const defaultMsg = "already running"

func (b BaseFailureTask) Action(tgBot *services.TgBot, update tgbotapi.Update) error {
	return tgBot.ReplyMessage(update, b.msg)
}

func NewCronTask(task services.Task, failureTask services.Task, sleepTime time.Duration) services.Task {
	return &CronTask{running: make(map[int64]struct{}), task: task, failureTask: failureTask, sleepTime: sleepTime}
}

func NewReplyCronAction(task services.Task, msg string, sleepTime time.Duration) services.Task {
	return NewCronTask(task, BaseFailureTask{msg: msg}, sleepTime)
}

func NewDefaultCronAction(task services.Task, sleepTime time.Duration) services.Task {
	return NewReplyCronAction(task, defaultMsg, sleepTime)
}

func (c *CronTask) Action(tgBot *services.TgBot, update tgbotapi.Update) error {
	_, ok := c.running[update.Message.Chat.ID]
	if !ok {
		c.running[update.Message.Chat.ID] = struct{}{}
		go func(running *map[int64]struct{}) {
			for {
				err := c.task.Action(tgBot, update)
				if err != nil {
					logrus.Errorf("[CronAction.Action] error: %v", err)
					c.m.Lock()
					delete(c.running, update.Message.Chat.ID)
					c.m.Unlock()
				}
				time.Sleep(c.sleepTime)
			}
		}(&c.running)
	} else {
		return c.failureTask.Action(tgBot, update)
	}
	return nil
}
