package notifier

import (
	"context"
	"errors"
	"fmt"

	"github.com/slack-go/slack"
)

type slackBroker struct {
	ctx context.Context
}

func newSlackBroker() *slackBroker {
	broker := &slackBroker{
		ctx: context.Background(),
	}
	return broker
}

func (m *slackBroker) isValid() bool {
	return true
}

func (m *slackBroker) sendNotification(n *Notification) error {
	if !m.isValid() {
		return errors.New("slack broker is not valid")
	}
	return slack.PostWebhookContext(m.ctx, n.Dest, &slack.WebhookMessage{
		Text: fmt.Sprintf("Threshold: %v\nCurrentSpend: %v\nUser: %v", n.Threshold, n.CurrentSpend, n.UserMail),
	})
}
