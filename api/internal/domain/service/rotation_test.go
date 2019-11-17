package service

import (
	"context"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"github.com/koind/banner-rotation/api/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRotationService_Add(t *testing.T) {
	rotationService := RotationService{
		RotationRepository: memory.NewRotationRepository(),
	}

	rotation := repository.Rotation{
		ID:          1,
		BannerID:    13,
		SlotID:      5,
		Description: "New rotation",
		CreateAt:    time.Now().UTC(),
	}

	newRotation, _ := rotationService.Add(context.Background(), rotation)
	assert.Equal(t, &rotation, newRotation)
}

func TestRotationService_SetTransition(t *testing.T) {
	statisticRepository := memory.NewStatisticRepository()
	rotationService := RotationService{
		RotationRepository: memory.NewRotationRepository(),
		StatisticService: &StatisticService{
			StatisticRepository: statisticRepository,
		},
		StatisticRepository: statisticRepository,
	}

	groupID := 8
	slotID := 5
	rotation := repository.Rotation{
		ID:          1,
		BannerID:    13,
		SlotID:      slotID,
		Description: "New rotation",
		CreateAt:    time.Now().UTC(),
	}

	expectedStatistic := repository.Statistic{
		ID:       1,
		Type:     repository.StatisticTypeClick,
		BannerID: 13,
		SlotID:   5,
		GroupID:  8,
	}

	err := rotationService.SetTransition(context.Background(), rotation, groupID)
	assert.Nil(t, err)

	statistics, _ := rotationService.StatisticRepository.FindAllBySlotIDAndGroupID(
		context.Background(),
		slotID,
		groupID,
	)

	expectedStatistic.CreateAt = statistics[0].CreateAt
	assert.Equal(t, &expectedStatistic, statistics[0])
}

func TestRotationService_Remove(t *testing.T) {
	rotationService := RotationService{
		RotationRepository: memory.NewRotationRepository(),
	}

	rotation := repository.Rotation{
		ID:          1,
		BannerID:    13,
		SlotID:      5,
		Description: "New rotation",
		CreateAt:    time.Now().UTC(),
	}

	newRotation, _ := rotationService.Add(context.Background(), rotation)
	assert.Equal(t, &rotation, newRotation)

	err := rotationService.Remove(context.Background(), rotation.BannerID)
	assert.Nil(t, err)
}

func TestRotationService_SelectBanner(t *testing.T) {
	statisticRepository := memory.StatisticRepository{
		DB: map[int]repository.Statistic{
			1: {
				ID:       1,
				Type:     repository.StatisticTypeView,
				BannerID: 1,
				SlotID:   1,
				GroupID:  1,
				CreateAt: time.Now().UTC(),
			},
			2: {
				ID:       2,
				Type:     repository.StatisticTypeClick,
				BannerID: 1,
				SlotID:   1,
				GroupID:  1,
				CreateAt: time.Now().UTC(),
			},
			3: {
				ID:       3,
				Type:     repository.StatisticTypeView,
				BannerID: 2,
				SlotID:   1,
				GroupID:  1,
				CreateAt: time.Now().UTC(),
			},
			4: {
				ID:       4,
				Type:     repository.StatisticTypeView,
				BannerID: 3,
				SlotID:   1,
				GroupID:  1,
				CreateAt: time.Now().UTC(),
			},
			5: {
				ID:       5,
				Type:     repository.StatisticTypeClick,
				BannerID: 3,
				SlotID:   1,
				GroupID:  1,
				CreateAt: time.Now().UTC(),
			},
			6: {
				ID:       6,
				Type:     repository.StatisticTypeView,
				BannerID: 1,
				SlotID:   1,
				GroupID:  1,
				CreateAt: time.Now().UTC(),
			},
			7: {
				ID:       7,
				Type:     repository.StatisticTypeView,
				BannerID: 3,
				SlotID:   1,
				GroupID:  1,
				CreateAt: time.Now().UTC(),
			},
			8: {
				ID:       8,
				Type:     repository.StatisticTypeClick,
				BannerID: 3,
				SlotID:   1,
				GroupID:  1,
				CreateAt: time.Now().UTC(),
			},
		},
		ID: 9,
	}

	rotationRepository := memory.RotationRepository{
		DB: map[int]repository.Rotation{
			1: {
				ID:          1,
				BannerID:    1,
				SlotID:      1,
				Description: "Banner 1",
				CreateAt:    time.Now().UTC(),
			},
			2: {
				ID:          2,
				BannerID:    2,
				SlotID:      1,
				Description: "Banner 2",
				CreateAt:    time.Now().UTC(),
			},
			3: {
				ID:          3,
				BannerID:    3,
				SlotID:      1,
				Description: "Banner 3",
				CreateAt:    time.Now().UTC(),
			},
		},
		ID: 4,
	}

	rotationService := RotationService{
		RotationRepository: &rotationRepository,
		StatisticService: &StatisticService{
			StatisticRepository: &statisticRepository,
		},
		StatisticRepository: &statisticRepository,
	}

	slotID := 1
	groupId := 1
	expectedBannerID := 3

	bannerID, _ := rotationService.SelectBanner(context.Background(), slotID, groupId)
	assert.Equal(t, expectedBannerID, bannerID, "")
}
