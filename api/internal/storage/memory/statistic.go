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
type StatisticRepository struct {
	sync.RWMutex
	DB map[int]repository.Statistic
	ID int
}

// Will return new memory statistics repository
func NewStatisticRepository() *StatisticRepository {
	return &StatisticRepository{
		DB: make(map[int]repository.Statistic),
		ID: 1,
	}
}

// Adds statistics
func (s *StatisticRepository) Add(
	ctx context.Context,
	statistic repository.Statistic,
) (*repository.Statistic, error) {
	s.Lock()
	defer s.Unlock()

	statistic.ID = s.ID
	s.DB[statistic.ID] = statistic
	s.ID++

	return &statistic, nil
}

// Find all the statistics by slot and group
func (s *StatisticRepository) FindAllBySlotIDAndGroupID(
	ctx context.Context,
	slotID int,
	groupID int,
) ([]*repository.Statistic, error) {
	s.RLock()
	defer s.RUnlock()

	if len(s.DB) <= 0 {
		return nil, nil
	}

	statistics := make([]*repository.Statistic, 0)

	for _, statistic := range s.DB {
		if statistic.SlotID == slotID && statistic.GroupID == groupID {
			statistic := statistic
			statistics = append(statistics, &statistic)
		}
	}

	return statistics, nil
}

// Removes statistics
func (s *StatisticRepository) Remove(ctx context.Context, ID int) error {
	s.Lock()
	defer s.Unlock()

	if _, has := s.DB[ID]; !has {
		return ErrStatisticNotFound
	}

	delete(s.DB, ID)

	return nil
}
