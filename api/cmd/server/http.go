package server

import (
	"context"
	"github.com/koind/banner-rotation/api/internal/config"
	"github.com/koind/banner-rotation/api/internal/db"
	"github.com/koind/banner-rotation/api/internal/domain/service"
	"github.com/koind/banner-rotation/api/internal/rabbit"
	"github.com/koind/banner-rotation/api/internal/storage/postgres"
	transport "github.com/koind/banner-rotation/api/internal/transport/http"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
	"time"
)

var HttpServerCmd = &cobra.Command{
	Use:   "http_server",
	Short: "Run http server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Init(config.Path)

		logger, err := zap.NewProduction()
		if err != nil {
			log.Fatal(err)
		}
		defer logger.Sync()

		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(cfg.Postgres.PingTimeout)*time.Millisecond,
		)
		defer cancel()

		pg, err := db.IntPostgres(ctx, config.Postgres(cfg.Postgres))
		if err != nil {
			log.Fatalf("failing to connect to the database %v", err)
		}

		conn, err := amqp.Dial(cfg.RabbitMQ.URL)
		if err != nil {
			log.Fatalf("failing to connect to the rabbitmq %v", err)
		}
		defer conn.Close()

		rotationRepository := postgres.NewRotationRepository(pg, *logger)
		statisticRepository := postgres.NewStatisticRepository(pg, *logger)
		statisticService := service.StatisticService{StatisticRepository: statisticRepository}
		rotationService := service.RotationService{
			StatisticService:    &statisticService,
			RotationRepository:  rotationRepository,
			StatisticRepository: statisticRepository,
		}

		publisher := rabbit.NewPublisher(*conn, cfg.RabbitMQ.ExchangeName, cfg.RabbitMQ.QueueName)
		httpRotationService := transport.NewHTTPRotationService(rotationService, publisher, logger)
		hs := transport.NewHTTPServer(httpRotationService, cfg.HTTPServer.GetDomain())

		logger.Error("Error starting app", zap.Error(hs.Start()))
	},
}

func init() {
	HttpServerCmd.Flags().StringVarP(
		&config.Path,
		"config",
		"c",
		"config/development/config.toml",
		"Path to toml configuration file",
	)
}
