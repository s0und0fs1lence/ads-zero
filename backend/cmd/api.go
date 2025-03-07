package cmd

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/s0und0fs1lence/ads-zero/pkg/api"
	"github.com/s0und0fs1lence/ads-zero/pkg/configuration"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{

	Use:   "start-api",
	Short: "Ads-Zero is a tool for fetching spend metrics from Advertising platform, and send alerts based on rules",

	RunE: func(cmd *cobra.Command, args []string) error {
		// Do Stuff Here

		conn, err := clickhouse.Open(&clickhouse.Options{
			Addr: []string{configuration.Config().GetString(configuration.ClickhouseHosts)}, //TODO: add proper params
			Auth: clickhouse.Auth{
				Database: configuration.Config().GetString(configuration.ClickhouseDb),
				Username: configuration.Config().GetString(configuration.ClickhouseUsername),
				Password: configuration.Config().GetString(configuration.ClickhousePassword),
			},

			Settings: clickhouse.Settings{
				"max_execution_time": 60,
			},
			Compression: &clickhouse.Compression{
				Method: clickhouse.CompressionZSTD,
				Level:  9,
			},
			DialTimeout:          time.Second * 30,
			MaxOpenConns:         30,
			MaxIdleConns:         5,
			ConnMaxLifetime:      time.Duration(10) * time.Minute,
			ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
			BlockBufferSize:      10,
			MaxCompressionBuffer: 10240,
			ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "ads-zero-backend", Version: "0.1"},
				},
			},
		})
		if err != nil {
			return err
		}
		if err := conn.Ping(context.Background()); err != nil {
			return err
		}

		db, err := db.NewClickhouseService(conn)
		if err != nil {
			return err
		}
		if err := api.StartAPI(db); err != nil {
			return err
		}
		return nil
	},
}

func init() {

	rootCmd.AddCommand(apiCmd)
}
