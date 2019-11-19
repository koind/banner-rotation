package server

import (
	"context"
	"github.com/koind/banner-rotation/api/internal/config"
	"github.com/koind/banner-rotation/api/internal/db"
	"github.com/koind/banner-rotation/api/internal/domain/service"
	"github.com/koind/banner-rotation/api/internal/storage/postgres"
	"github.com/koind/banner-rotation/api/internal/transport/grpc"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"time"
)

var GrpcServerCmd = &cobra.Command{
	Use:   "grpc_server",
	Short: "Run grpc server",
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

		rotationRepository := postgres.NewRotationRepository(pg, *logger)
		statisticRepository := postgres.NewStatisticRepository(pg, *logger)
		statisticService := service.StatisticService{StatisticRepository: statisticRepository}
		rotationService := service.RotationService{
			StatisticService:    &statisticService,
			RotationRepository:  rotationRepository,
			StatisticRepository: statisticRepository,
		}

		gs := grpc.RotationServer{
			Domain:          cfg.GRPCServer.GetDomain(),
			RotationService: rotationService,
		}

		err = gs.Start()
		if err != nil {
			log.Fatalf("failed to serve %v", err)
		}
	},
}

func init() {
	GrpcServerCmd.Flags().StringVarP(
		&config.Path,
		"config",
		"c",
		"config/development/config.toml",
		"Path to toml configuration file",
	)
}
