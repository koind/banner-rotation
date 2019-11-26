package server

import (
	"context"
	"github.com/koind/banner-rotation/api/internal/config"
	"github.com/koind/banner-rotation/api/internal/db"
	"github.com/koind/banner-rotation/api/internal/domain/service"
	"github.com/koind/banner-rotation/api/internal/rabbit"
	"github.com/koind/banner-rotation/api/internal/storage/postgres"
	"github.com/koind/banner-rotation/api/internal/transport/grpc"
	"github.com/koind/banner-rotation/api/internal/transport/http"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

// Declaring commands to start server
var RunServerCmd = &cobra.Command{
	Use:   "run",
	Short: "Run server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Init(config.Path)
		rotationService, publisher, logger := Init(cfg)
		serverType := os.Getenv("SERVER_TYPE")

		switch serverType {
		case "HTTP":
			httpRotationService := http.NewHTTPRotationService(*rotationService, publisher, logger)
			hs := http.NewHTTPServer(httpRotationService, cfg.HTTPServer.GetDomain())

			logger.Error("Error starting http server", zap.Error(hs.Start()))
		case "GRPC":
			gs := grpc.NewGRPCServer(cfg.GRPCServer.GetDomain(), *rotationService, publisher, logger)

			logger.Error("Error starting grpc server", zap.Error(gs.Start()))
		default:
			log.Fatal("Specified the wrong server type")
		}
	},
}

// Returns the initialized objects needed to start the server
func Init(cfg config.Options) (*service.RotationService, *rabbit.Publisher, *zap.Logger) {
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
	statisticsRepository := postgres.NewStatisticsRepository(pg, *logger)
	statisticsService := service.StatisticsService{StatisticsRepository: statisticsRepository}
	publisher := rabbit.NewPublisher(conn, cfg.RabbitMQ.ExchangeName, cfg.RabbitMQ.QueueName)
	rotationService := service.RotationService{
		StatisticsService:    &statisticsService,
		RotationRepository:   rotationRepository,
		StatisticsRepository: statisticsRepository,
	}

	return &rotationService, publisher, logger
}

// When initializing parse the path to the configuration
func init() {
	RunServerCmd.Flags().StringVarP(
		&config.Path,
		"config",
		"c",
		"config/development/config.toml",
		"Path to toml configuration file",
	)
}
