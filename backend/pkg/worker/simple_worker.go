package worker

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/rs/zerolog/log"
	"github.com/s0und0fs1lence/ads-zero/pkg/common"
	"github.com/s0und0fs1lence/ads-zero/pkg/configuration"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
	"github.com/s0und0fs1lence/ads-zero/pkg/fetcher"
	"github.com/s0und0fs1lence/ads-zero/pkg/notifier"
)

type SimpleWorker interface {
	Run()
	Done() chan struct{}
}

type simpleWorker struct {
	workerID      string
	wg            sync.WaitGroup
	done          chan struct{}
	ctx           context.Context
	cancel        context.CancelFunc
	dbSvc         db.DbService
	messageBroker notifier.MessageBroker
	ticker        *time.Ticker
	queryBuilder  *sqlbuilder.SelectBuilder
}

func (k *simpleWorker) teardown() error {
	k.wg.Wait()
	log.Info().Msg("done all the goroutines... exiting now")
	return nil
}

func (k *simpleWorker) processRuleExecution(client fetcher.Client, task *common.FetchTask) error {

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

func (k *simpleWorker) fetchData(msg common.ScheduleMessage) {
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
	if task == nil {
		log.Error().Any("message", msg).Msg("task is nil")
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

func (k *simpleWorker) scheduleFetches() error {
	defer k.wg.Done()
	clients, err := k.dbSvc.GetClientsByQuery(k.queryBuilder)
	if err != nil {
		log.Error().Err(err).Msg("failed to get clients")
		return err
	}

	dt := time.Now().UTC()
	for _, c := range clients {

		msg := common.ScheduleMessage{
			ClientID: c.ClientID,
			Start:    dt,
			End:      dt,
		}
		k.wg.Add(1)
		go k.fetchData(msg)

	}
	return nil
}

// ConsumeClaim implements Worker.
func (k *simpleWorker) Run() {
	for {
		if err := k.ctx.Err(); err != nil {
			if err := k.teardown(); err != nil {
				log.Error().Err(err).Msg("got error while awaiting the schedules")
			}
			return
		}
		select {
		case <-k.ctx.Done():
			if err := k.teardown(); err != nil {
				log.Error().Err(err).Msg("got error while awaiting the schedules")
			}

			log.Warn().Msg("context ended. exiting now...")
			return

		case <-k.done:
			if err := k.teardown(); err != nil {
				log.Error().Err(err).Msg("got error while awaiting the schedules")
			}

			log.Warn().Msg("scheduler sigterm received. exiting now...")
			return
		case tm := <-k.ticker.C:
			log.Info().Time("tick", tm).Msg("starting scheduler..")
			k.wg.Add(1)
			go k.scheduleFetches()

		}
	}
}

// Done implements Scheduler.
func (k *simpleWorker) Done() chan struct{} {
	return k.done
}

func NewSimpleWorker(ctx context.Context, cancel context.CancelFunc, svc db.DbService) SimpleWorker {
	qb := sqlbuilder.NewSelectBuilder().Select("*").From("clients")
	qb.Where(
		qb.EQ("deleted", false),
	)
	broker, err := notifier.NewMessageBroker()
	if err != nil {
		log.Error().Err(err).Msg("could not create the message broker")
		return nil
	}
	return &simpleWorker{
		workerID:      uuid.NewString(),
		done:          make(chan struct{}),
		ctx:           ctx,
		cancel:        cancel,
		dbSvc:         svc,
		ticker:        time.NewTicker(time.Minute * time.Duration(configuration.Config().GetInt(configuration.SimpleWorkerTickInterval))),
		queryBuilder:  qb,
		messageBroker: broker,
	}
}
