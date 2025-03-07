package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/s0und0fs1lence/ads-zero/pkg/common"
	"github.com/s0und0fs1lence/ads-zero/pkg/configuration"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
	"github.com/wneessen/go-mail"
)

func TestUnmarshall(t *testing.T) {

	var s common.ScheduleMessage
	bts := []byte(`
		{
			"client_id": "test",
			"start" : "2023-01-01",
			"end": "2023-01-02"
		}
	`)
	if err := json.Unmarshal(bts, &s); err != nil {
		t.Fatal(err)
	}

}

func TestMarshall(t *testing.T) {

	s := common.ScheduleMessage{
		ClientID: "TEST",
		Start:    time.Now().UTC(),
		End:      time.Now().UTC(),
	}

	if bts, err := json.Marshal(&s); err != nil {
		t.Fatal(err)
	} else {
		var testStruct common.ScheduleMessage
		if err := json.Unmarshal(bts, &testStruct); err != nil {
			t.Fatal(err)
		}
		if s.ClientID != testStruct.ClientID {
			t.Fatal("client id differ")
		}
		if s.Start.Format(time.DateOnly) != testStruct.Start.Format(time.DateOnly) {
			t.Fatal("start differ")
		}
		if s.End.Format(time.DateOnly) != testStruct.End.Format(time.DateOnly) {
			t.Fatal("end differ")
		}
	}

}

func TestMailClient(t *testing.T) {
	client, err := mail.NewClient(
		configuration.Config().GetString(configuration.MailHost),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(configuration.Config().GetString(configuration.MailUsername)),
		mail.WithPassword(configuration.Config().GetString(configuration.MailPassword)),
	)
	if err != nil {
		t.Fatal(err)
	}
	message := mail.NewMsg()
	if err := message.From("toni@tester.com"); err != nil {
		t.Fatalf("failed to set FROM address: %s", err)
	}
	if err := message.To("test@gmail.com"); err != nil {
		t.Fatalf("failed to set TO address: %s", err)
	}
	f, err := os.ReadFile("/workspaces/ads-zero/pkg/worker/template.html")
	if err != nil {
		t.Fatal(err)
	}
	d := templateData{
		Threshold:    "aaa",
		CurrentSpend: 50,
		User:         "test",
	}

	htmlTpl, err := template.New("htmltpl").Parse(string(f))
	if err != nil {
		t.Fatalf("failed to parse text template: %s", err)
	}
	if err := message.AddAlternativeHTMLTemplate(htmlTpl, d); err != nil {
		t.Fatal(err)
	}
	message.Subject("This is my first test mail with go-mail!")
	// message.SetBodyString(mail.TypeTextPlain, "This will be the content of the mail.")
	if err := client.DialAndSend(message); err != nil {
		t.Fatalf("failed to deliver mail: %s", err)
	}
	fmt.Printf("client: %v\n", client)
}

func TestSimpleWorker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	conn, err := db.NewClickhouseService(nil)
	if err != nil {
		t.Fatal(err)
	}
	worker := NewSimpleWorker(ctx, cancel, conn)
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)
	go func() {
		worker.Run()
	}()
	log.Info().Msg("started the consumer...")

	<-cancelChan
	log.Warn().Msg("signal received... sending termination to the worker")
	worker.Done() <- struct{}{}
	cancel()
	<-worker.Done()
	log.Info().Msg("exiting now...")
}
