package service

import (
	"context"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"github.com/pkg/errors"
	"time"
)

// The service interface statistics
type StatisticServiceInterface interface {
	// Saves the statistics
	Save(ctx context.Context, rotation repository.Rotation, groupID int, statisticType int) (*repository.Statistic, error)
}

// Statistic service
type StatisticService struct {
	StatisticRepository repository.StatisticRepositoryInterface
}

// Saves the statistics
func (s *StatisticService) Save(
	ctx context.Context,
	rotation repository.Rotation,
	groupID int,
	statisticType int,
) (*repository.Statistic, error) {
	statistic := repository.Statistic{
		Type:     statisticType,
		BannerID: rotation.BannerID,
		SlotID:   rotation.SlotID,
		GroupID:  groupID,
		CreateAt: time.Now().UTC(),
	}

	newStatistic, err := s.StatisticRepository.Add(ctx, statistic)
	if err != nil {
		return nil, errors.Wrap(err, "error saving statistics")
	}

	return newStatistic, nil
}
