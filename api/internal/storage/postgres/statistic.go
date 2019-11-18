package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	queryInsertStatistic = `INSERT INTO statistics(type, banner_id, slot_id, group_id, create_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	queryFindAllBySlotIDAndGroupID = `SELECT * FROM statistics WHERE slot_id=$1 AND group_id=$2`
	queryRemoveByStatisticID       = `DELETE FROM statistics WHERE id=$1`
)

// Postgres statistic repository
type StatisticRepository struct {
	DB     *sqlx.DB
	logger zap.Logger
}

// Returns the postgres statistic repository
func NewStatisticRepository(db *sqlx.DB, logger zap.Logger) *StatisticRepository {
	return &StatisticRepository{
		DB:     db,
		logger: logger,
	}
}

// Adds statistics
func (s *StatisticRepository) Add(ctx context.Context, statistic repository.Statistic) (*repository.Statistic, error) {
	if ctx.Err() == context.Canceled {
		s.logger.Info(
			"Adding a statistics was canceled due to context cancellation",
			zap.Int("BannerID", statistic.BannerID),
		)

		return nil, errors.New("adding a statistics was canceled due to context cancellation")
	}

	err := s.DB.QueryRowContext(
		ctx,
		queryInsertStatistic,
		statistic.Type,
		statistic.BannerID,
		statistic.SlotID,
		statistic.GroupID,
		statistic.CreateAt,
	).Scan(&statistic.ID)
	if err != nil {
		return nil, errors.Wrap(err, "error when adding statistics")
	}

	return &statistic, nil
}

// Find all the statistics by slot and group
func (s *StatisticRepository) FindAllBySlotIDAndGroupID(
	ctx context.Context,
	slotID int,
	groupID int,
) ([]*repository.Statistic, error) {
	if ctx.Err() == context.Canceled {
		s.logger.Info(
			"Search for all statistics was interrupted due to context cancellation",
			zap.Int("slotID", slotID),
			zap.Int("groupID", groupID),
		)

		return nil, errors.New("search for all statistics was interrupted due to context cancellation")
	}

	rows, err := s.DB.QueryxContext(ctx, queryFindAllBySlotIDAndGroupID, slotID, groupID)
	if err != nil {
		return nil, errors.Wrap(err, "error when searching statistics by slotId and groupId")
	}

	statistics := make([]*repository.Statistic, 0)

	for rows.Next() {
		var statistic repository.Statistic
		err := rows.StructScan(&statistic)
		if err != nil {
			return nil, errors.Wrap(err, "error while scanning results to structure")
		}

		statistics = append(statistics, &statistic)
	}

	return statistics, nil
}

// Removes statistics
func (s *StatisticRepository) Remove(ctx context.Context, ID int) error {
	if ctx.Err() == context.Canceled {
		s.logger.Info(
			"Removal statistic was interrupted due to the cancellation context",
			zap.Int("ID", ID),
		)

		return errors.New("removal statistic was interrupted due to the cancellation context")
	}

	_, err := s.DB.ExecContext(ctx, queryRemoveByStatisticID, ID)
	if err != nil {
		return errors.Wrap(err, "error when remove statistic")
	}

	return nil
}
