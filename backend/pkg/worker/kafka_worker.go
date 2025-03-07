package worker

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"github.com/s0und0fs1lence/ads-zero/pkg/common"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
	"github.com/s0und0fs1lence/ads-zero/pkg/fetcher"
	"github.com/s0und0fs1lence/ads-zero/pkg/notifier"
)

type Worker interface {
	Setup(session sarama.ConsumerGroupSession) error
	Cleanup(session sarama.ConsumerGroupSession) error
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
	Ready() chan bool
	ReinitChan()
	EnsureFlush() error
}

type templateData struct {
	Threshold    any
	CurrentSpend float64
	User         string
}

type kafkaWorker struct {
	workerID      string
	wg            sync.WaitGroup
	ready         chan bool
	ctx           context.Context
	cancel        context.CancelFunc
	dbSvc         db.DbService
	producer      sarama.SyncProducer
	messageBroker notifier.MessageBroker
}

func (k *kafkaWorker) teardown() error {
	k.wg.Wait()
	log.Info().Msg("done all the goroutines... exiting now")
	return nil
}

// Cleanup implements Worker.
func (k *kafkaWorker) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Info().Msg("cleanup was called...")

	return k.teardown()

}

func (k *kafkaWorker) processRuleExecution(client fetcher.Client, task *common.FetchTask) error {

	results, err := client.ExecuteRules(task)
	if err != nil {
		return err
	}
	log.Info().Any("total_spend", task.GetTotalSpend()).Any("rule_result", results).Msg("processing rule execution")
	notified := false
	for _, result := range results {
		if result.Result {
			channel, val := client.GetNotificationChannel()
			log.Info().Str("notification_channel", channel).Str("value", val).Msg("sending notification")
			k.messageBroker.SendNotification(&notifier.Notification{
				CurrentSpend: task.GetTotalSpend(),
				Threshold:    result.Threshold,
				RuleName:     result.RuleName,
				RuleID:       result.RuleID,
				DestType:     channel,
				Dest:         val,
			})
			notified = true

		}
		if notified {
			break
		}
	}
	return nil

}

func (k *kafkaWorker) fetchData(msg common.ScheduleMessage) {
	defer k.wg.Done()
	client := fetcher.NewClient().WithDbSvc(k.dbSvc).WithUserId(msg.ClientID)
	if client == nil {
		log.Error().Any("message", msg).Msg("could not create the client")
		return
	}
	if !client.IsValid() {
		log.Error().Any("message", msg).Err(client.GetError()).Msg("")
		return
	}
	task, err := client.FetchData(msg.Start, msg.End)
	if err != nil {
		log.Error().Any("message", msg).Err(err).Msg("")
		return

	}

	if err := client.SaveAccountData(task.Accounts); err != nil {
		log.Error().Any("message", msg).Err(err).Msg("")
		return
	}
	if err := client.SaveCampaignData(task.Campaigns); err != nil {
		log.Error().Any("message", msg).Err(err).Msg("")
		return
	}
	if err := k.processRuleExecution(client, task); err != nil {
		log.Error().Any("message", msg).Err(err).Msg("")
		return
	}
	log.Info().Any("message", msg).Msg("done fetching data")

}

// ConsumeClaim implements Worker.
func (k *kafkaWorker) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		if session.Context().Err() != nil {
			log.Error().Err(session.Context().Err()).Msg("the session has been closed")
			return nil
		}
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Error().Msg("message channel was closed")
				return k.teardown()
			}
			var msg common.ScheduleMessage
			if err := json.Unmarshal(message.Value, &msg); err != nil {
				log.Error().Err(err).Int64("offset", message.Offset).Str("topic", message.Topic).Msg("could not decode the message topic")
				// the message is not valid... we should just skip the processing...
				continue
			}
			log.Info().Interface("message", msg).Msg("starting fetch for the given message...")
			k.wg.Add(1)
			go k.fetchData(msg)
		case <-session.Context().Done():
			return k.teardown()
		case <-k.ctx.Done():
			return k.teardown()

		}
	}
}

// EnsureFlush implements Worker.
func (k *kafkaWorker) EnsureFlush() error {
	panic("unimplemented")
}

// Ready implements Worker.
func (k *kafkaWorker) Ready() chan bool {
	return k.ready
}

// ReinitChan implements Worker.
func (k *kafkaWorker) ReinitChan() {
	k.ready = make(chan bool)
}

// Setup implements Worker.
func (k *kafkaWorker) Setup(session sarama.ConsumerGroupSession) error {
	k.workerID = session.MemberID()
	broker, err := notifier.NewMessageBroker()
	if err != nil {
		return err
	}
	k.messageBroker = broker
	close(k.ready)
	k.wg = sync.WaitGroup{}
	log.Info().Str("worker_id", k.workerID).Msg("worker initialized")
	return nil
}

func NewKafkaWorker(ctx context.Context, cancel context.CancelFunc, producer sarama.SyncProducer, svc db.DbService) Worker {
	return &kafkaWorker{
		workerID: "test",
		ready:    make(chan bool),
		ctx:      ctx,
		cancel:   cancel,
		dbSvc:    svc,
		producer: producer,
	}
}
