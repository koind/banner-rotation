package repository

import (
	"context"
	"time"
)

// Banner model
type Banner struct {
	ID      int
	GroupID int
	Views   int
	Clicks  float64
}

// Rotation model
type Rotation struct {
	ID          int       `json:"id" db:"id"`
	BannerID    int       `json:"bannerId" db:"banner_id"`
	SlotID      int       `json:"slotId" db:"slot_id"`
	Description string    `json:"description" db:"description"`
	CreateAt    time.Time `json:"createAt" db:"create_at"`
}

// The repository interface rotation
type RotationRepositoryInterface interface {
	// Adds a new banner to the rotation in this slot
	Add(ctx context.Context, rotation Rotation) (*Rotation, error)

	// Find all rotations by slot id
	FindAllBySlotID(ctx context.Context, slotID int) ([]*Rotation, error)

	// Removes the banner from the rotation
	Remove(ctx context.Context, bannerID int) error
}
