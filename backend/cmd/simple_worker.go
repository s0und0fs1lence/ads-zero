package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/s0und0fs1lence/ads-zero/pkg/configuration"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
	"github.com/s0und0fs1lence/ads-zero/pkg/worker"
	"github.com/spf13/cobra"
)

var simpleWorkerCmd = &cobra.Command{

	Use:   "start-simple-worker",
	Short: "Ads-Zero is a tool for fetching spend metrics from Advertising platform, and send alerts based on rules",

	RunE: func(cmd *cobra.Command, args []string) error {
		// Do Stuff Here

		level, err := zerolog.ParseLevel(configuration.Config().GetString(configuration.LogLevel))
		if err != nil {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		} else {
			zerolog.SetGlobalLevel(level)
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancelChan := make(chan os.Signal, 1)
		signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)
		db, err := db.NewClickhouseService(nil)
		if err != nil {
			cancel()
			return err
		}

		wrk := worker.NewSimpleWorker(ctx, cancel, db)

		go func() {
			wrk.Run()
		}()
		log.Info().Msg("worker started...")
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

	rootCmd.AddCommand(simpleWorkerCmd)
}
