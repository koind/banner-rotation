package service

import (
	"context"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"github.com/koind/banner-rotation/api/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStatisticsService_Save(t *testing.T) {
	statisticsService := StatisticsService{
		StatisticsRepository: memory.NewStatisticsRepository(),
	}

	testCases := map[string]struct {
		rotation           repository.Rotation
		groupID            int
		statisticsType     int
		expectedStatistics repository.Statistics
	}{
		"Save statistics with type type view": {
			rotation: repository.Rotation{
				ID:          1,
				BannerID:    13,
				SlotID:      5,
				Description: "New rotation",
				CreatedAt:   time.Now().UTC(),
			},
			groupID:        2,
			statisticsType: repository.StatisticsTypeView,
			expectedStatistics: repository.Statistics{
				ID:       1,
				Type:     repository.StatisticsTypeView,
				BannerID: 13,
				SlotID:   5,
				GroupID:  2,
			},
		},
		"Save statistics with type type click": {
			rotation: repository.Rotation{
				ID:          1,
				BannerID:    13,
				SlotID:      5,
				Description: "New rotation",
				CreatedAt:   time.Now().UTC(),
			},
			groupID:        2,
			statisticsType: repository.StatisticsTypeClick,
			expectedStatistics: repository.Statistics{
				ID:       2,
				Type:     repository.StatisticsTypeClick,
				BannerID: 13,
				SlotID:   5,
				GroupID:  2,
			},
		},
	}

	for _, testCase := range testCases {
		statistics, _ := statisticsService.Save(
			context.Background(),
			testCase.rotation,
			testCase.groupID,
			testCase.statisticsType,
		)

		testCase.expectedStatistics.CreatedAt = statistics.CreatedAt

		assert.Equal(t, &testCase.expectedStatistics, statistics, "values must match")
	}
}
