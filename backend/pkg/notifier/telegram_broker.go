package notifier

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/s0und0fs1lence/ads-zero/pkg/configuration"
)

type telegramBroker struct {
	client *bot.Bot
	ctx    context.Context
}

func newTelegramBroker() (*telegramBroker, error) {
	bot, err := bot.New(configuration.Config().GetString(configuration.TelegramBotToken))
	if err != nil {
		return nil, err
	}
	broker := &telegramBroker{
		client: bot,
		ctx:    context.Background(),
	}
	return broker, nil
}

func (m *telegramBroker) isValid() bool {
	return m.client != nil
}

func (m *telegramBroker) sendNotification(n *Notification) error {
	if !m.isValid() {
		return errors.New("mail broker is not valid")
	}
	//TODO: Implement telegram notification

	msg, err := m.client.SendMessage(m.ctx, &bot.SendMessageParams{
		Text:   fmt.Sprintf("Threshold: %v\nCurrentSpend: %v\nUser: %v", n.Threshold, n.CurrentSpend, n.UserMail),
		ChatID: n.Dest,
	})
	if err != nil {
		return err
	}
	fmt.Printf("msg: %v\n", msg)

	return nil
}
