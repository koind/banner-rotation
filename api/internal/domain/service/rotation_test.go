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

	statistic, err := rotationService.SetTransition(context.Background(), rotation, groupID)
	assert.Nil(t, err)

	expectedStatistic.CreateAt = statistic.CreateAt
	assert.Equal(t, &expectedStatistic, statistic)
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
	var testCases = []struct {
		rotationRepository  memory.RotationRepository
		statisticRepository memory.StatisticRepository
		slotID              int
		groupID             int
		err                 error
		expectedBannerID    int
	}{
		{
			rotationRepository:  memory.RotationRepository{},
			statisticRepository: memory.StatisticRepository{},
			slotID:              1,
			groupID:             1,
			err:                 ErrRotationsListEmpty,
			expectedBannerID:    0,
		},
		{
			rotationRepository: memory.RotationRepository{
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
				},
				ID: 3,
			},
			statisticRepository: memory.StatisticRepository{
				DB: map[int]repository.Statistic{
					1: {
						ID:       1,
						Type:     repository.StatisticTypeView,
						BannerID: 1,
						SlotID:   1,
						GroupID:  1,
						CreateAt: time.Now().UTC(),
					},
				},
			},
			slotID:           1,
			groupID:          1,
			err:              nil,
			expectedBannerID: 2,
		},
		{
			rotationRepository: memory.RotationRepository{
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
			},
			statisticRepository: memory.StatisticRepository{
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
			},
			slotID:           1,
			groupID:          1,
			err:              nil,
			expectedBannerID: 3,
		},
		{
			rotationRepository: memory.RotationRepository{
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
						SlotID:      2,
						Description: "Banner 2",
						CreateAt:    time.Now().UTC(),
					},
					3: {
						ID:          3,
						BannerID:    3,
						SlotID:      3,
						Description: "Banner 3",
						CreateAt:    time.Now().UTC(),
					},
					4: {
						ID:          4,
						BannerID:    4,
						SlotID:      2,
						Description: "Banner 4",
						CreateAt:    time.Now().UTC(),
					},
				},
				ID: 5,
			},
			statisticRepository: memory.StatisticRepository{
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
						SlotID:   2,
						GroupID:  4,
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
						BannerID: 4,
						SlotID:   2,
						GroupID:  4,
						CreateAt: time.Now().UTC(),
					},
					8: {
						ID:       8,
						Type:     repository.StatisticTypeClick,
						BannerID: 4,
						SlotID:   2,
						GroupID:  4,
						CreateAt: time.Now().UTC(),
					},
				},
				ID: 9,
			},
			slotID:           2,
			groupID:          4,
			err:              nil,
			expectedBannerID: 4,
		},
	}

	for _, testCase := range testCases {
		rotationService := RotationService{
			RotationRepository: &testCase.rotationRepository,
			StatisticService: &StatisticService{
				StatisticRepository: &testCase.statisticRepository,
			},
			StatisticRepository: &testCase.statisticRepository,
		}

		bannerID, _, err := rotationService.SelectBanner(context.Background(), testCase.slotID, testCase.groupID)

		if err != nil {
			assert.Error(t, testCase.err, &err)
		} else {
			assert.Equal(t, testCase.expectedBannerID, bannerID, "banners ids must match")
		}
	}
}
