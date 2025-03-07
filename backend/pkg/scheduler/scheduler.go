package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/rs/zerolog/log"
	"github.com/s0und0fs1lence/ads-zero/pkg/common"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
	"golang.org/x/sync/errgroup"
)

type Scheduler interface {
	Run()
	Done() chan struct{}
}

type kafkaScheduler struct {
	workerID     string
	wg           *errgroup.Group
	done         chan struct{}
	ctx          context.Context
	cancel       context.CancelFunc
	dbSvc        db.DbService
	producer     sarama.AsyncProducer
	topic        string
	ticker       *time.Ticker
	queryBuilder *sqlbuilder.SelectBuilder
}

// Done implements Scheduler.
func (k *kafkaScheduler) Done() chan struct{} {
	return k.done
}

func NewKafkaScheduler(ctx context.Context, cancel context.CancelFunc, producer sarama.AsyncProducer, svc db.DbService, topic string) Scheduler {
	qb := sqlbuilder.NewSelectBuilder().Select("*").From("clients")
	qb.Where(
		qb.EQ("deleted", false),
		// qb.EQ("stripe_subscription_status","ACTIVE"),
	)

	return &kafkaScheduler{
		workerID:     uuid.NewString(),
		wg:           new(errgroup.Group),
		done:         make(chan struct{}),
		ctx:          ctx,
		cancel:       cancel,
		dbSvc:        svc,
		producer:     producer,
		ticker:       time.NewTicker(time.Minute * 15),
		queryBuilder: qb,
		topic:        topic,
	}
}

// Run implements Scheduler.
func (k *kafkaScheduler) Run() {
	defer close(k.done)
	for {
		if err := k.ctx.Err(); err != nil {
			if err := k.wg.Wait(); err != nil {
				log.Error().Err(err).Msg("got error while awaiting the schedules..")
			}
			log.Error().Err(err).Msg("got context error... exiting")
			return
		}
		select {
		case <-k.ctx.Done():
			if err := k.wg.Wait(); err != nil {
				log.Error().Err(err).Msg("got error while awaiting the schedules..")
			}
			log.Warn().Msg("context ended. exiting now...")
			return

		case <-k.done:
			if err := k.wg.Wait(); err != nil {
				log.Error().Err(err).Msg("got error while awaiting the schedules..")
			}
			log.Warn().Msg("scheduler sigterm received. exiting now...")
			return
		case tm := <-k.ticker.C:
			log.Info().Time("tick", tm).Msg("starting scheduler..")

			k.wg.Go(k.scheduleFetches)

			// if err := k.scheduleFetches(); err != nil {
			// 	log.Error().Err(err).Msg("could not schedule the fetches...")
			// }
			// log.Info().Msg("done scheduling messages...")
		}

	}
}

func (k *kafkaScheduler) scheduleFetches() error {
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
		bts, err := json.Marshal(&msg)
		if err != nil {

			return err
		}
		k.producer.Input() <- &sarama.ProducerMessage{
			Topic: k.topic,
			Value: sarama.ByteEncoder(bts),
		}
	}
	return nil
}
