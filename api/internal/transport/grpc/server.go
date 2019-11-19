package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"github.com/koind/banner-rotation/api/internal/domain/service"
	pb "github.com/koind/banner-rotation/api/internal/transport/grpc/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

// GRPC rotation service
type RotationServer struct {
	Domain          string
	RotationService service.RotationService
}

// Adds a banner in the rotation
func (s *RotationServer) AddBanner(ctx context.Context, req *pb.RotationRequest) (*pb.RotationResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, errors.New("client cancelled, abandoning.")
	}

	rotation := repository.Rotation{
		BannerID:    int(req.GetBannerId()),
		SlotID:      int(req.GetSlotId()),
		Description: req.GetDescription(),
	}

	rotation.SetDatetimeOfCreate()

	newRotation, err := s.RotationService.Add(ctx, rotation)
	if err != nil {
		return nil, err
	}

	createAt, err := ptypes.TimestampProto(newRotation.CreateAt)
	if err != nil {
		return nil, err
	}

	rotationResp := &pb.RotationResponse{
		Id:          int32(newRotation.ID),
		BannerId:    int32(newRotation.BannerID),
		SlotId:      int32(newRotation.SlotID),
		Description: newRotation.Description,
		CreateAt:    createAt,
	}

	return rotationResp, nil
}

// Sets the transition on the banner
func (s *RotationServer) SetTransition(ctx context.Context, t *pb.Transition) (*pb.Status, error) {
	if ctx.Err() == context.Canceled {
		return nil, errors.New("client cancelled, abandoning.")
	}

	bannerID := int(t.GetBannerId())
	groupID := int(t.GetGroupId())

	rotation, err := s.RotationService.RotationRepository.FindOneByBannerID(ctx, bannerID)
	if err != nil {
		return nil, err
	}

	err = s.RotationService.SetTransition(ctx, *rotation, groupID)
	if err != nil {
		return nil, err
	}

	return &pb.Status{Status: "ok"}, nil
}

// Selects a banner to display
func (s *RotationServer) SelectBanner(ctx context.Context, sl *pb.Select) (*pb.Banner, error) {
	if ctx.Err() == context.Canceled {
		return nil, errors.New("client cancelled, abandoning.")
	}

	slotID := int(sl.GetSlotId())
	groupID := int(sl.GetGroupId())

	bannerID, err := s.RotationService.SelectBanner(ctx, slotID, groupID)
	if err != nil {
		return nil, err
	}

	banner := &pb.Banner{
		Id: int32(bannerID),
	}

	return banner, nil
}

// Removes the banner from the rotation
func (s *RotationServer) RemoveBanner(ctx context.Context, banner *pb.Banner) (*pb.Status, error) {
	if ctx.Err() == context.Canceled {
		return nil, errors.New("client cancelled, abandoning.")
	}

	bannerID := int(banner.GetId())

	err := s.RotationService.Remove(ctx, bannerID)
	if err != nil {
		return nil, err
	}

	return &pb.Status{Status: "ok"}, nil
}

// Start fires up the grpc server
func (s *RotationServer) Start() error {
	gs := grpc.NewServer()
	reflection.Register(gs)

	l, err := net.Listen("tcp", s.Domain)
	if err != nil {
		return err
	}

	pb.RegisterRotationServer(gs, s)

	return gs.Serve(l)
}
