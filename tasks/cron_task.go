package tasks

import (
	"regexp"
	"sync"
	"time"

	"github.com/apala4i/simple_tg_bot_service/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type cronTask struct {
	running     map[int64]struct{}
	m           sync.Mutex
	task        services.Task
	failureTask services.Task
	stopTask    services.Task
	startTask   services.Task
	sleepTime   time.Duration
}

type baseFailureTask struct {
	msg string
}

const (
	defaultFailureMsg = "already running!\n"
	defaultStartMsg   = "started!\n"
	defaultStoppedMsg = "stopped!\n"
)

func newBaseFailureTask(msg string) services.Task {
	return NewTask(func(tgBot *services.TgBot, update tgbotapi.Update) error {
		return tgBot.ReplyMessage(update, msg)
	}, baseName, "")
}

func newBaseStopTask(taskName string, msg string) services.Task {
	name := "/stop_" + taskName[1:]
	return NewTask(func(tgBot *services.TgBot, update tgbotapi.Update) error {
		return tgBot.ReplyMessage(update, msg)
	}, name, "stop "+taskName)
}

func newBaseStartTask(msg string) services.Task {
	return NewTask(func(tgBot *services.TgBot, update tgbotapi.Update) error {
		return tgBot.ReplyMessage(update, msg)
	}, baseName, "")
}

func NewCronTask(task services.Task,
	failureTask services.Task,
	stopTask services.Task,
	startTask services.Task,
	sleepTime time.Duration) services.Task {
	return &cronTask{
		running:     make(map[int64]struct{}),
		task:        task,
		failureTask: failureTask,
		startTask:   startTask,
		stopTask:    stopTask,
		sleepTime:   sleepTime,
	}
}

func NewReplyCronAction(task services.Task,
	failureMsg string,
	startMsg string,
	stopMsg string,
	sleepTime time.Duration) services.Task {
	return NewCronTask(task,
		newBaseFailureTask(failureMsg),
		newBaseStopTask(task.GetNamePattern().String(), stopMsg),
		newBaseStartTask(startMsg), sleepTime)
}

// noinspection GoUnusedExportedFunction
func NewDefaultCronTask(task services.Task, sleepTime time.Duration) services.Task {
	return NewReplyCronAction(task, defaultFailureMsg, defaultStartMsg, defaultStoppedMsg, sleepTime)
}

func (c *cronTask) Action(tgBot *services.TgBot, update tgbotapi.Update) error {
	_, ok := c.running[update.Message.Chat.ID]
	if ok && c.stopTask.CompareName(update.Message.Text) {
		c.deleteFromQueue(update.Message.Chat.ID)
		err := c.stopTask.Action(tgBot, update)
		if err != nil {
			logrus.Errorf("stop action failed. error: %v", err)
		}
	} else if !ok {
		if err := c.startTask.Action(tgBot, update); err != nil {
			logrus.Errorf("start action failed. err: %v", err)
		}
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

func (c *cronTask) deleteFromQueue(chatId int64) {
	c.m.Lock()
	delete(c.running, chatId)
	c.m.Unlock()
}

func (c *cronTask) GetNamePattern() *regexp.Regexp {
	return c.task.GetNamePattern()
}

func (c *cronTask) GetDescription() string {
	return c.task.GetDescription() + "\n" + c.stopTask.GetDescription()
}

func (c *cronTask) CompareName(name string) bool {
	return c.task.GetNamePattern().Match([]byte(name)) || c.stopTask.GetNamePattern().Match([]byte(name))
}
