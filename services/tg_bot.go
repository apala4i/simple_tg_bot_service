package services

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	Api *tgbotapi.BotAPI
	m   sync.Mutex
}

func NewTgBot(bot *tgbotapi.BotAPI) *TgBot {
	return &TgBot{
		Api: bot,
		m:   sync.Mutex{},
	}
}

func (c *TgBot) SendMessage(chatId int64, msgText string) error {
	c.lock()
	defer func() {
		c.unlock()
	}()

	_, err := c.Api.Send(tgbotapi.NewMessage(chatId, msgText))
	if err != nil {
		return err
	}
	return nil
}

func (c *TgBot) ReplyMessage(update tgbotapi.Update, msgText string) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyToMessageID = update.Message.MessageID
	c.lock()
	defer func() {
		c.unlock()
	}()
	_, err := c.Api.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *TgBot) GetBot() *tgbotapi.BotAPI {
	return c.Api
}

func (c *TgBot) lock() {
	c.m.Lock()
}

func (c *TgBot) unlock() {
	c.m.Unlock()
}
