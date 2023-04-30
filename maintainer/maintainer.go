package maintainer

import "github.com/apala4i/simple_tg_bot_service/services"

type Maintainer interface {
	AddServer(botName string, bot *services.TgBotServer)
	DeleteServer(botName string)
	Start()
}
type maintainer struct {
	tgBots map[string]*services.TgBotServer
}

func NewMaintainer() Maintainer {
	return &maintainer{tgBots: make(map[string]*services.TgBotServer)}
}

func (m *maintainer) AddServer(botName string, bot *services.TgBotServer) {
	m.tgBots[botName] = bot
}

func (m *maintainer) DeleteServer(botName string) {
	delete(m.tgBots, botName)
}

func (m *maintainer) Start() {
	for _, bot := range m.tgBots {
		bot.Start()
	}
}
