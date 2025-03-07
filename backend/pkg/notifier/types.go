package notifier

type Notification struct {
	Subject      string
	UserMail     string
	CurrentSpend float64
	Threshold    any
	RuleName     string
	RuleID       string
	DestType     string
	// Expected values:
	// - if DestType is "mail", then Dest is the email address
	// - if DestType is "slack", then Dest is the slack webhook
	// - if DestType is "telegram", then Dest is the telegram chat id
	Dest string
}

type MessageBroker interface {
	SendNotification(*Notification) error
}
