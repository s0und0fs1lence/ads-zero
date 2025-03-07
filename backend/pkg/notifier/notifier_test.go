package notifier

import "testing"

func TestTelegram(t *testing.T) {
	broker, err := NewMessageBroker()
	if err != nil {
		t.Errorf("Error creating message broker: %v", err)
	}
	if broker == nil {
		t.Errorf("Broker is nil")
	}
	broker.SendNotification(&Notification{
		UserMail:     "",
		CurrentSpend: 0,
		Threshold:    0,
		RuleName:     "",
		RuleID:       "",
		DestType:     "slack",
		Dest:         "https://hooks.slack.com/services/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",
	})
}
