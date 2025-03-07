package cmd

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/s0und0fs1lence/ads-zero/pkg/configuration"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
	"github.com/s0und0fs1lence/ads-zero/pkg/scheduler"
	"github.com/spf13/cobra"
)

var schedulerCmd = &cobra.Command{

	Use:   "start-scheduler",
	Short: "Ads-Zero is a tool for fetching spend metrics from Advertising platform, and send alerts based on rules",

	RunE: func(cmd *cobra.Command, args []string) error {
		// Do Stuff Here

		level, err := zerolog.ParseLevel(configuration.Config().GetString(configuration.LogLevel))
		if err != nil {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		} else {
			zerolog.SetGlobalLevel(level)
		}

		topic := configuration.Config().GetString(configuration.KafkaTopic)

		conf := sarama.NewConfig()
		conf.Net.SASL.Mechanism = sarama.SASLTypePlaintext //"PLAIN"
		conf.Net.SASL.Enable = true
		conf.Net.SASL.User = configuration.Config().GetString(configuration.KafkaUsername)
		conf.Net.SASL.Password = configuration.Config().GetString(configuration.KafkaPassword)
		conf.Producer.Compression = sarama.CompressionZSTD
		conf.Producer.CompressionLevel = 9
		// conf.Producer.Transaction.ID = uuid.NewString()
		// conf.Producer.Idempotent = true
		// conf.Producer.RequiredAcks = sarama.WaitForAll
		conf.Metadata.Full = true

		ctx, cancel := context.WithCancel(context.Background())

		cancelChan := make(chan os.Signal, 1)
		signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)

		db, err := db.NewClickhouseService(nil)
		if err != nil {
			cancel()
			return err
		}
		kafkaBrokers := strings.Split(configuration.Config().GetString(configuration.KafkaUrl), ",")
		producer, err := sarama.NewAsyncProducer(kafkaBrokers, conf)

		if err != nil {
			cancel()
			return err
		}
		producer.Input()

		wrk := scheduler.NewKafkaScheduler(ctx, cancel, producer, db, topic)

		go func() {
			wrk.Run()
		}()
		log.Info().Msg("started the consumer...")

		<-cancelChan
		log.Warn().Msg("signal received... sending termination to the worker")
		wrk.Done() <- struct{}{}
		cancel()
		<-wrk.Done()
		log.Info().Msg("exiting now...")
		return nil
	},
}

func init() {

	rootCmd.AddCommand(schedulerCmd)
}
