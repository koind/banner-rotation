package memory

import (
	"context"
	"errors"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"sync"
)

var (
	ErrStatisticNotFound = errors.New("statistic not found")
)

// Memory statistics repository
type StatisticsRepository struct {
	sync.RWMutex
	DB map[int]repository.Statistics
	ID int
}

// Will return new memory statistics repository
func NewStatisticsRepository() *StatisticsRepository {
	return &StatisticsRepository{
		DB: make(map[int]repository.Statistics),
		ID: 1,
	}
}

// Adds statistics
func (s *StatisticsRepository) Add(
	ctx context.Context,
	statistics repository.Statistics,
) (*repository.Statistics, error) {
	s.Lock()
	defer s.Unlock()

	statistics.ID = s.ID
	s.DB[statistics.ID] = statistics
	s.ID++

	return &statistics, nil
}

// Find all the statistics by slot and group
func (s *StatisticsRepository) FindAllBySlotIDAndGroupID(
	ctx context.Context,
	slotID int,
	groupID int,
) ([]*repository.Statistics, error) {
	s.RLock()
	defer s.RUnlock()

	if len(s.DB) <= 0 {
		return nil, nil
	}

	statisticsList := make([]*repository.Statistics, 0)

	for _, statistics := range s.DB {
		if statistics.SlotID == slotID && statistics.GroupID == groupID {
			statistics := statistics
			statisticsList = append(statisticsList, &statistics)
		}
	}

	return statisticsList, nil
}

// Removes statistics
func (s *StatisticsRepository) Remove(ctx context.Context, ID int) error {
	s.Lock()
	defer s.Unlock()

	if _, has := s.DB[ID]; !has {
		return ErrStatisticNotFound
	}

	delete(s.DB, ID)

	return nil
}
