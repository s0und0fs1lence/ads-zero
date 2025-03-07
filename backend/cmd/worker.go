package cmd

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/s0und0fs1lence/ads-zero/pkg/configuration"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
	"github.com/s0und0fs1lence/ads-zero/pkg/worker"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{

	Use:   "start-worker",
	Short: "Ads-Zero is a tool for fetching spend metrics from Advertising platform, and send alerts based on rules",

	RunE: func(cmd *cobra.Command, args []string) error {
		// Do Stuff Here

		level, err := zerolog.ParseLevel(configuration.Config().GetString(configuration.LogLevel))
		if err != nil {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		} else {
			zerolog.SetGlobalLevel(level)
		}

		conf := sarama.NewConfig()
		conf.Net.SASL.Mechanism = sarama.SASLTypePlaintext //"PLAIN"
		conf.Net.SASL.Enable = true
		conf.Net.SASL.User = configuration.Config().GetString(configuration.KafkaUsername)
		conf.Net.SASL.Password = configuration.Config().GetString(configuration.KafkaPassword)
		conf.Producer.Compression = sarama.CompressionZSTD
		conf.Producer.CompressionLevel = 9
		conf.Producer.RequiredAcks = sarama.WaitForLocal
		conf.Producer.Return.Successes = true
		conf.Producer.Return.Errors = true
		conf.Metadata.Full = true
		conf.Consumer.Offsets.Initial = sarama.OffsetNewest
		//TODO: add manual commit strategy, to have a retry mechanism if the fetch had any error...
		conf.Consumer.Offsets.AutoCommit = struct {
			Enable   bool
			Interval time.Duration
		}{Enable: true, Interval: time.Second * 10}

		ctx, cancel := context.WithCancel(context.Background())
		cancelChan := make(chan os.Signal, 1)
		signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)
		db, err := db.NewClickhouseService(nil)
		if err != nil {
			cancel()
			return err
		}
		kafkaBrokers := strings.Split(configuration.Config().GetString(configuration.KafkaUrl), ",")
		consumerGroupName := configuration.Config().GetString(configuration.KafkaConsumerGroupName)
		kafkaConsumerGroup, err := sarama.NewConsumerGroup(kafkaBrokers, consumerGroupName, conf)
		if err != nil {
			cancel()
			return err
		}

		producer, err := sarama.NewSyncProducer(kafkaBrokers, conf)
		if err != nil {
			cancel()
			return err
		}

		wrk := worker.NewKafkaWorker(ctx, cancel, producer, db)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if err := kafkaConsumerGroup.Consume(ctx, []string{configuration.Config().GetString(configuration.KafkaTopic)}, wrk); err != nil {
					if errors.Is(err, sarama.ErrClosedConsumerGroup) {
						return
					}
					log.Error().Err(err).Msg("GOT ERROR")

				}
				if ctx.Err() != nil {
					return
				}
				wrk.ReinitChan()
			}
		}()
		<-wrk.Ready()
		log.Info().Msg("started the consumer...")
		keepRunning := true
		for keepRunning {
			if ctx.Err() != nil {
				log.Error().Err(ctx.Err()).Msg("terminating: context error")
				if err := kafkaConsumerGroup.Close(); err != nil {
					log.Error().Err(err).Msg("could not close kafka consumer group")
					return err
				}
				return ctx.Err()

			}
			select {
			case <-ctx.Done():
				log.Warn().Msg("termiting: context cancelled")
				keepRunning = false
			case <-cancelChan:
				log.Warn().Msg("termiting: received signal")
				keepRunning = false
			}
		}
		cancel()
		wg.Wait()
		if err := kafkaConsumerGroup.Close(); err != nil {
			log.Error().Err(err).Msg("could not close kafka consumer group")
			return err
		}
		log.Info().Msg("exiting now...")
		return nil
	},
}

func init() {

	rootCmd.AddCommand(workerCmd)
}
