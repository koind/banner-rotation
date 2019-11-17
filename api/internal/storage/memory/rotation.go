package memory

import (
	"context"
	"errors"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"sync"
)

var (
	ErrRotationNotFound = errors.New("rotation not found")
)

// Memory rotation repository
type RotationRepository struct {
	sync.RWMutex
	DB map[int]repository.Rotation
	ID int
}

// Will return new memory rotation repository
func NewRotationRepository() *RotationRepository {
	return &RotationRepository{
		DB: make(map[int]repository.Rotation),
		ID: 1,
	}
}

// Adds a new banner to the rotation in this slot
func (r *RotationRepository) Add(ctx context.Context, rotation repository.Rotation) (*repository.Rotation, error) {
	r.Lock()
	defer r.Unlock()

	rotation.ID = r.ID
	r.DB[rotation.ID] = rotation
	r.ID++

	return &rotation, nil
}

// Find all rotations by slot id
func (r *RotationRepository) FindAllBySlotID(ctx context.Context, slotID int) ([]*repository.Rotation, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.DB) <= 0 {
		return nil, nil
	}

	rotations := make([]*repository.Rotation, 0)

	for _, rotation := range r.DB {
		if rotation.SlotID == slotID {
			rotation := rotation
			rotations = append(rotations, &rotation)
		}
	}

	return rotations, nil
}

// Removes the banner from the rotation
func (r *RotationRepository) Remove(ctx context.Context, bannerID int) error {
	r.Lock()
	defer r.Unlock()

	rotationIDS := make([]int, 0)

	for _, rotation := range r.DB {
		if rotation.BannerID == bannerID {
			rotationIDS = append(rotationIDS, rotation.ID)
		}
	}

	if len(rotationIDS) <= 0 {
		return ErrRotationNotFound
	}

	for _, rotationID := range rotationIDS {
		delete(r.DB, rotationID)
	}

	return nil
}
