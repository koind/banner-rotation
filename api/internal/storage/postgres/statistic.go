package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	queryInsertStatistic = `INSERT INTO statistics(type, banner_id, slot_id, group_id, created_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	queryFindAllBySlotIDAndGroupID = `SELECT * FROM statistics WHERE slot_id=$1 AND group_id=$2`
	queryRemoveByStatisticID       = `DELETE FROM statistics WHERE id=$1`
)

// Postgres statistics repository
type StatisticsRepository struct {
	DB     *sqlx.DB
	logger zap.Logger
}

// Returns the postgres statistics repository
func NewStatisticsRepository(db *sqlx.DB, logger zap.Logger) *StatisticsRepository {
	return &StatisticsRepository{
		DB:     db,
		logger: logger,
	}
}

// Adds statistics
func (s *StatisticsRepository) Add(ctx context.Context, statistics repository.Statistics) (*repository.Statistics, error) {
	if ctx.Err() == context.Canceled {
		s.logger.Info(
			"Adding a statistics was canceled due to context cancellation",
			zap.Int("BannerID", statistics.BannerID),
		)

		return nil, errors.New("adding a statistics was canceled due to context cancellation")
	}

	err := s.DB.QueryRowContext(
		ctx,
		queryInsertStatistic,
		statistics.Type,
		statistics.BannerID,
		statistics.SlotID,
		statistics.GroupID,
		statistics.CreatedAt,
	).Scan(&statistics.ID)
	if err != nil {
		return nil, errors.Wrap(err, "error when adding statistics")
	}

	return &statistics, nil
}

// Find all the statistics by slot and group
func (s *StatisticsRepository) FindAllBySlotIDAndGroupID(
	ctx context.Context,
	slotID int,
	groupID int,
) ([]*repository.Statistics, error) {
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

	statisticsList := make([]*repository.Statistics, 0)

	for rows.Next() {
		var statistics repository.Statistics
		err := rows.StructScan(&statistics)
		if err != nil {
			return nil, errors.Wrap(err, "error while scanning results to structure")
		}

		statisticsList = append(statisticsList, &statistics)
	}

	return statisticsList, nil
}

// Removes statistics
func (s *StatisticsRepository) Remove(ctx context.Context, ID int) error {
	if ctx.Err() == context.Canceled {
		s.logger.Info(
			"Removal statistics was interrupted due to the cancellation context",
			zap.Int("ID", ID),
		)

		return errors.New("removal statistics was interrupted due to the cancellation context")
	}

	_, err := s.DB.ExecContext(ctx, queryRemoveByStatisticID, ID)
	if err != nil {
		return errors.Wrap(err, "error when remove statistics")
	}

	return nil
}
