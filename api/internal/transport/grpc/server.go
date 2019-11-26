package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"github.com/koind/banner-rotation/api/internal/domain/service"
	"github.com/koind/banner-rotation/api/internal/rabbit"
	pb "github.com/koind/banner-rotation/api/internal/transport/grpc/api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

// GRPC rotation service
type GrpcServer struct {
	domain          string
	rotationService service.RotationService
	publisher       rabbit.PublisherInterface
	logger          *zap.Logger
}

// NewGRPCServer returns grpc server that wraps rotation business logic
func NewGRPCServer(
	domain string,
	rotationService service.RotationService,
	publisher rabbit.PublisherInterface,
	logger *zap.Logger,
) *GrpcServer {
	return &GrpcServer{
		domain:          domain,
		rotationService: rotationService,
		publisher:       publisher,
		logger:          logger,
	}
}

// Adds a banner in the rotation
func (s *GrpcServer) AddBanner(ctx context.Context, req *pb.RotationRequest) (*pb.RotationResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, errors.New("client cancelled, abandoning.")
	}

	rotation := repository.Rotation{
		BannerID:    int(req.GetBannerId()),
		SlotID:      int(req.GetSlotId()),
		Description: req.GetDescription(),
	}

	rotation.SetDatetimeOfCreate()

	newRotation, err := s.rotationService.Add(ctx, rotation)
	if err != nil {
		return nil, err
	}

	createdAt, err := ptypes.TimestampProto(newRotation.CreatedAt)
	if err != nil {
		return nil, err
	}

	rotationResp := &pb.RotationResponse{
		Id:          int32(newRotation.ID),
		BannerId:    int32(newRotation.BannerID),
		SlotId:      int32(newRotation.SlotID),
		Description: newRotation.Description,
		CreateAt:    createdAt,
	}

	return rotationResp, nil
}

// Sets the transition on the banner
func (s *GrpcServer) SetTransition(ctx context.Context, t *pb.Transition) (*pb.Status, error) {
	if ctx.Err() == context.Canceled {
		return nil, errors.New("client cancelled, abandoning.")
	}

	bannerID := int(t.GetBannerId())
	groupID := int(t.GetGroupId())

	rotation, err := s.rotationService.RotationRepository.FindOneByBannerID(ctx, bannerID)
	if err != nil {
		return nil, err
	}

	statistics, err := s.rotationService.SetTransition(ctx, *rotation, groupID)
	if err != nil {
		return nil, err
	}

	err = s.publisher.Publish(ctx, *statistics)
	if err != nil {
		s.logger.Error(
			"Failed to send message to queue",
			zap.Error(err),
		)
	}

	return &pb.Status{Status: "ok"}, nil
}

// Selects a banner to display
func (s *GrpcServer) SelectBanner(ctx context.Context, sl *pb.Select) (*pb.Banner, error) {
	if ctx.Err() == context.Canceled {
		return nil, errors.New("client cancelled, abandoning.")
	}

	slotID := int(sl.GetSlotId())
	groupID := int(sl.GetGroupId())

	bannerID, statistics, err := s.rotationService.SelectBanner(ctx, slotID, groupID)
	if err != nil {
		return nil, err
	}

	err = s.publisher.Publish(ctx, *statistics)
	if err != nil {
		s.logger.Error(
			"Failed to send message to queue",
			zap.Error(err),
		)
	}

	banner := &pb.Banner{
		Id: int32(bannerID),
	}

	return banner, nil
}

// Removes the banner from the rotation
func (s *GrpcServer) RemoveBanner(ctx context.Context, banner *pb.Banner) (*pb.Status, error) {
	if ctx.Err() == context.Canceled {
		return nil, errors.New("client cancelled, abandoning.")
	}

	bannerID := int(banner.GetId())

	err := s.rotationService.Remove(ctx, bannerID)
	if err != nil {
		return nil, err
	}

	return &pb.Status{Status: "ok"}, nil
}

// Start fires up the grpc server
func (s *GrpcServer) Start() error {
	gs := grpc.NewServer()
	reflection.Register(gs)

	l, err := net.Listen("tcp", s.domain)
	if err != nil {
		return err
	}

	pb.RegisterRotationServer(gs, s)

	return gs.Serve(l)
}
