package repository

import (
	"context"
	"time"
)

const (
	// Type of statistics view
	StatisticTypeView = 1

	// Type of statistics click
	StatisticTypeClick = 2
)

// Statistic model
type Statistic struct {
	ID        int       `json:"id" db:"id"`
	Type      int       `json:"type" db:"type"`
	BannerID  int       `json:"bannerId" db:"banner_id"`
	SlotID    int       `json:"slotId" db:"slot_id"`
	GroupID   int       `json:"groupId" db:"group_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// Is the view type
func (s *Statistic) IsTypeView() bool {
	return s.Type == StatisticTypeView
}

// Is the click type
func (s *Statistic) IsTypeClick() bool {
	return s.Type == StatisticTypeClick
}

// The repository interface statistics
type StatisticRepositoryInterface interface {
	// Adds statistics
	Add(ctx context.Context, statistic Statistic) (*Statistic, error)

	// Find all the statistics by slot and group
	FindAllBySlotIDAndGroupID(ctx context.Context, slotID int, groupID int) ([]*Statistic, error)

	// Removes statistics
	Remove(ctx context.Context, ID int) error
}
