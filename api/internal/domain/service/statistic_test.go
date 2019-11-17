package service

import (
	"context"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"github.com/koind/banner-rotation/api/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStatisticService_Save(t *testing.T) {
	statisticService := StatisticService{
		StatisticRepository: memory.NewStatisticRepository(),
	}

	testCases := map[string]struct {
		rotation          repository.Rotation
		groupID           int
		statisticType     int
		expectedStatistic repository.Statistic
	}{
		"Save statistic with type type view": {
			rotation: repository.Rotation{
				ID:          1,
				BannerID:    13,
				SlotID:      5,
				Description: "New rotation",
				CreateAt:    time.Now().UTC(),
			},
			groupID:       2,
			statisticType: repository.StatisticTypeView,
			expectedStatistic: repository.Statistic{
				ID:       1,
				Type:     repository.StatisticTypeView,
				BannerID: 13,
				SlotID:   5,
				GroupID:  2,
			},
		},
		"Save statistic with type type click": {
			rotation: repository.Rotation{
				ID:          1,
				BannerID:    13,
				SlotID:      5,
				Description: "New rotation",
				CreateAt:    time.Now().UTC(),
			},
			groupID:       2,
			statisticType: repository.StatisticTypeClick,
			expectedStatistic: repository.Statistic{
				ID:       2,
				Type:     repository.StatisticTypeClick,
				BannerID: 13,
				SlotID:   5,
				GroupID:  2,
			},
		},
	}

	for _, testCase := range testCases {
		statistic, _ := statisticService.Save(
			context.Background(),
			testCase.rotation,
			testCase.groupID,
			testCase.statisticType,
		)

		testCase.expectedStatistic.CreateAt = statistic.CreateAt

		assert.Equal(t, &testCase.expectedStatistic, statistic, "values must match")
	}
}
