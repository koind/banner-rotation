package service

import (
	"context"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"github.com/pkg/errors"
	"time"
)

// The service interface statistics
type StatisticsServiceInterface interface {
	// Saves the statistics
	Save(ctx context.Context, rotation repository.Rotation, groupID int, statisticType int) (*repository.Statistics, error)
}

// Statistics service
type StatisticsService struct {
	StatisticsRepository repository.StatisticsRepositoryInterface
}

// Saves the statistics
func (s *StatisticsService) Save(
	ctx context.Context,
	rotation repository.Rotation,
	groupID int,
	statisticType int,
) (*repository.Statistics, error) {
	statistics := repository.Statistics{
		Type:      statisticType,
		BannerID:  rotation.BannerID,
		SlotID:    rotation.SlotID,
		GroupID:   groupID,
		CreatedAt: time.Now().UTC(),
	}

	newStatistics, err := s.StatisticsRepository.Add(ctx, statistics)
	if err != nil {
		return nil, errors.Wrap(err, "error saving statistics")
	}

	return newStatistics, nil
}
