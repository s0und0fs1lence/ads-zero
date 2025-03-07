package notifier

import (
	"errors"
	"html/template"
	"os"

	"github.com/s0und0fs1lence/ads-zero/pkg/configuration"
	"github.com/wneessen/go-mail"
)

type mailBroker struct {
	mailClient   *mail.Client
	htmlTemplate *template.Template
	fromEmail    string
}

type mailTemplateData struct {
	Threshold    any
	CurrentSpend float64
	User         string
}

func newMailBroker() (*mailBroker, error) {
	client, err := mail.NewClient(
		configuration.Config().GetString(configuration.MailHost),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(configuration.Config().GetString(configuration.MailUsername)),
		mail.WithPassword(configuration.Config().GetString(configuration.MailPassword)),
	)
	if err != nil {
		return nil, err
	}
	broker := &mailBroker{
		mailClient: client,
	}
	bts, err := os.ReadFile(configuration.Config().GetString(configuration.MailTemplate))
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("email").Parse(string(bts))
	if err != nil {
		return nil, err
	}
	broker.htmlTemplate = tmpl
	broker.fromEmail = "no-reply@adszero.com"
	return broker, nil
}

func (m *mailBroker) isValid() bool {
	if m.mailClient == nil {
		return false
	}
	if m.htmlTemplate == nil {
		return false
	}
	return true
}

func (m *mailBroker) sendNotification(n *Notification) error {
	if !m.isValid() {
		return errors.New("mail broker is not valid")
	}
	data := mailTemplateData{
		Threshold:    n.Threshold,
		CurrentSpend: n.CurrentSpend,
		User:         n.UserMail,
	}
	message := mail.NewMsg()
	if err := message.From(m.fromEmail); err != nil {
		return err
	}
	if err := message.To(n.UserMail); err != nil {
		return err
	}
	if err := message.AddAlternativeHTMLTemplate(m.htmlTemplate, data); err != nil {
		return err
	}
	if n.Subject == "" {
		n.Subject = "Alert: You have reached your spending threshold"
	}
	message.Subject(n.Subject)
	// message.SetBodyString(mail.TypeTextPlain, "This will be the content of the mail.")
	if err := m.mailClient.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
