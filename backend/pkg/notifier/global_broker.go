package notifier

import "errors"

type globalBroker struct {
	mailBroker  *mailBroker
	tgBroker    *telegramBroker
	slackBroker *slackBroker
}

func NewMessageBroker() (MessageBroker, error) {
	mb := &globalBroker{}
	mail, _ := newMailBroker()
	mb.mailBroker = mail
	tg, _ := newTelegramBroker()
	mb.tgBroker = tg
	mb.slackBroker = newSlackBroker()
	return mb, nil
}

// SendNotification implements MessageBroker.
func (g *globalBroker) SendNotification(n *Notification) error {
	if n == nil {
		return errors.New("notification is nil")
	}

	switch n.DestType {
	case "mail":
		if g.mailBroker == nil {
			return errors.New("mail broker is not initialized")
		}
		return g.mailBroker.sendNotification(n)
	case "telegram":
		if g.tgBroker == nil {
			return errors.New("telegram broker is not initialized")
		}
		return g.tgBroker.sendNotification(n)
	case "slack":
		if g.slackBroker == nil {
			return errors.New("slack broker is not initialized")
		}
		return g.slackBroker.sendNotification(n)
	default:
		return errors.New("unknown destination type")
	}
}
