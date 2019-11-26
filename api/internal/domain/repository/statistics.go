package repository

import (
	"context"
	"time"
)

const (
	// Type of statistics view
	StatisticsTypeView = 1

	// Type of statistics click
	StatisticsTypeClick = 2
)

// Statistics model
type Statistics struct {
	ID        int       `json:"id" db:"id"`
	Type      int       `json:"type" db:"type"`
	BannerID  int       `json:"bannerId" db:"banner_id"`
	SlotID    int       `json:"slotId" db:"slot_id"`
	GroupID   int       `json:"groupId" db:"group_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// Is the view type
func (s *Statistics) IsTypeView() bool {
	return s.Type == StatisticsTypeView
}

// Is the click type
func (s *Statistics) IsTypeClick() bool {
	return s.Type == StatisticsTypeClick
}

// The repository interface statistics
type StatisticsRepositoryInterface interface {
	// Adds statistics
	Add(ctx context.Context, statistics Statistics) (*Statistics, error)

	// Find all the statistics by slot and group
	FindAllBySlotIDAndGroupID(ctx context.Context, slotID int, groupID int) ([]*Statistics, error)

	// Removes statistics
	Remove(ctx context.Context, ID int) error
}
