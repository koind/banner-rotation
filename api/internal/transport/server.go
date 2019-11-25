package transport

import (
	"context"
	"github.com/koind/banner-rotation/api/internal/config"
	"github.com/koind/banner-rotation/api/internal/db"
	"github.com/koind/banner-rotation/api/internal/domain/service"
	"github.com/koind/banner-rotation/api/internal/rabbit"
	"github.com/koind/banner-rotation/api/internal/storage/postgres"
	"github.com/koind/banner-rotation/api/internal/transport/grpc"
	"github.com/koind/banner-rotation/api/internal/transport/http"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
	"time"
)

// Server struct
type Server struct {
	rotationService *service.RotationService
	publisher       *rabbit.Publisher
	logger          *zap.Logger
	config          config.Options
}

// Returns new server
func NewServer(cfg config.Options) *Server {
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
	publisher := rabbit.NewPublisher(*conn, cfg.RabbitMQ.ExchangeName, cfg.RabbitMQ.QueueName)
	rotationService := service.RotationService{
		StatisticService:    &statisticService,
		RotationRepository:  rotationRepository,
		StatisticRepository: statisticRepository,
	}

	return &Server{
		rotationService: &rotationService,
		publisher:       publisher,
		logger:          logger,
	}
}

// Run specified server
func (s *Server) Run(serverType string) {
	switch serverType {
	case "HTTP":
		httpRotationService := http.NewHTTPRotationService(*s.rotationService, s.publisher, s.logger)
		hs := http.NewHTTPServer(httpRotationService, s.config.HTTPServer.GetDomain())

		s.logger.Error("Error starting http server", zap.Error(hs.Start()))
	case "GRPC":
		gs := grpc.NewGRPCServer(s.config.GRPCServer.GetDomain(), *s.rotationService, s.publisher, s.logger)

		s.logger.Error("Error starting grpc server", zap.Error(gs.Start()))
	default:
		log.Fatal("Specified the wrong server type")
	}
}
