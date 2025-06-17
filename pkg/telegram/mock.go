package telegram

import tgbotapi "github.com/OvyFlash/telegram-bot-api"

type Mock struct {
	updates chan tgbotapi.Update
}

func (m *Mock) Send(_ tgbotapi.Chattable) (tgbotapi.Message, error) {
	return tgbotapi.Message{}, nil
}

func (m *Mock) GetUpdatesChan(_ tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	m.updates = make(chan tgbotapi.Update)

	return m.updates
}

func (m *Mock) SendUpdates(update tgbotapi.Update) {
	m.updates <- update
}

func (m *Mock) CloseUpdates() {
	close(m.updates)
}

func (m *Mock) Self() tgbotapi.User {
	return tgbotapi.User{
		ID:        1,
		IsBot:     false,
		FirstName: "Mock",
		UserName:  "Mock",
	}
}
