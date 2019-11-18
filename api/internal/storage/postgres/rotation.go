package postgres

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	queryInsertRotation = `INSERT INTO rotations(banner_id, slot_id, description, create_at)
		VALUES ($1, $2, $3, $4) RETURNING id`
	queryFindRotationByBannerID = `SELECT * FROM rotations WHERE banner_id=$1`
	queryFindAllBySlotID        = `SELECT * FROM rotations WHERE slot_id=$1`
	queryRemoveByBannerID       = `DELETE FROM rotations WHERE banner_id=$1`
)

// Postgres rotation repository
type RotationRepository struct {
	DB     *sqlx.DB
	logger zap.Logger
}

// Returns the postgres rotation repository
func NewRotationRepository(db *sqlx.DB, logger zap.Logger) *RotationRepository {
	return &RotationRepository{
		DB:     db,
		logger: logger,
	}
}

// Adds a new banner to the rotation in this slot
func (r *RotationRepository) Add(ctx context.Context, rotation repository.Rotation) (*repository.Rotation, error) {
	if ctx.Err() == context.Canceled {
		r.logger.Info(
			"Adding a banner to the rotation was canceled due to context cancellation",
			zap.Int("BannerID", rotation.BannerID),
		)

		return nil, errors.New("adding a banner to the rotation was canceled due to context cancellation")
	}

	err := r.DB.QueryRowContext(
		ctx,
		queryInsertRotation,
		rotation.BannerID,
		rotation.SlotID,
		rotation.Description,
		rotation.CreateAt,
	).Scan(&rotation.ID)
	if err != nil {
		return nil, errors.Wrap(err, "error when adding banner in the rotation")
	}

	return &rotation, nil
}

// Find one rotation by banner id
func (r *RotationRepository) FindOneByBannerID(ctx context.Context, bannerID int) (*repository.Rotation, error) {
	if ctx.Err() == context.Canceled {
		r.logger.Info(
			"Find one rotation was interrupted due to context cancellation",
			zap.Int("bannerID", bannerID),
		)

		return nil, errors.New("find one rotation was interrupted due to context cancellation")
	}

	row := r.DB.QueryRowContext(ctx, queryFindRotationByBannerID, bannerID)

	rotation := new(repository.Rotation)
	err := row.Scan(
		&rotation.ID,
		&rotation.BannerID,
		&rotation.SlotID,
		&rotation.Description,
		&rotation.CreateAt,
	)

	if err == sql.ErrNoRows {
		r.logger.Warn(
			"Could not find rotation by bannerID",
			zap.Int("bannerID", bannerID),
		)

		return nil, errors.Wrap(err, "could not find rotation by bannerID")
	} else if err != nil {
		r.logger.Warn(
			"Error when searching for rotation by bannerID",
			zap.Error(err),
			zap.Int("bannerID", bannerID),
		)

		return nil, errors.Wrap(err, "error when searching for rotation by bannerID")
	}

	return rotation, nil
}

// Find all rotations by slot id
func (r *RotationRepository) FindAllBySlotID(ctx context.Context, slotID int) ([]*repository.Rotation, error) {
	if ctx.Err() == context.Canceled {
		r.logger.Info(
			"Search for all banners was interrupted due to context cancellation",
			zap.Int("slotID", slotID),
		)

		return nil, errors.New("search for all banners was interrupted due to context cancellation")
	}

	rows, err := r.DB.QueryxContext(ctx, queryFindAllBySlotID, slotID)
	if err != nil {
		return nil, errors.Wrap(err, "error when searching for rotations by slotId")
	}

	rotations := make([]*repository.Rotation, 0)

	for rows.Next() {
		var rotation repository.Rotation
		err := rows.StructScan(&rotation)
		if err != nil {
			return nil, errors.Wrap(err, "error while scanning results to structure")
		}

		rotations = append(rotations, &rotation)
	}

	return rotations, nil
}

// Removes the banner from the rotation
func (r *RotationRepository) Remove(ctx context.Context, bannerID int) error {
	if ctx.Err() == context.Canceled {
		r.logger.Info(
			"Removal rotation of a banner was interrupted due to the cancellation context",
			zap.Int("bannerID", bannerID),
		)

		return errors.New("removal rotation of a banner was interrupted due to the cancellation context")
	}

	_, err := r.DB.ExecContext(ctx, queryRemoveByBannerID, bannerID)
	if err != nil {
		return errors.Wrap(err, "error when remove banner rotation")
	}

	return nil
}
