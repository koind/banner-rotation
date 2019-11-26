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
		CreatedAt:   time.Now().UTC(),
	}

	newRotation, _ := rotationService.Add(context.Background(), rotation)
	assert.Equal(t, &rotation, newRotation)
}

func TestRotationService_SetTransition(t *testing.T) {
	statisticsRepository := memory.NewStatisticsRepository()
	rotationService := RotationService{
		RotationRepository: memory.NewRotationRepository(),
		StatisticsService: &StatisticsService{
			StatisticsRepository: statisticsRepository,
		},
		StatisticsRepository: statisticsRepository,
	}

	groupID := 8
	slotID := 5
	rotation := repository.Rotation{
		ID:          1,
		BannerID:    13,
		SlotID:      slotID,
		Description: "New rotation",
		CreatedAt:   time.Now().UTC(),
	}

	expectedStatistics := repository.Statistics{
		ID:       1,
		Type:     repository.StatisticsTypeClick,
		BannerID: 13,
		SlotID:   5,
		GroupID:  8,
	}

	statistics, err := rotationService.SetTransition(context.Background(), rotation, groupID)
	assert.Nil(t, err)

	expectedStatistics.CreatedAt = statistics.CreatedAt
	assert.Equal(t, &expectedStatistics, statistics)
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
		CreatedAt:   time.Now().UTC(),
	}

	newRotation, _ := rotationService.Add(context.Background(), rotation)
	assert.Equal(t, &rotation, newRotation)

	err := rotationService.Remove(context.Background(), rotation.BannerID)
	assert.Nil(t, err)
}

func TestRotationService_SelectBanner(t *testing.T) {
	var testCases = []struct {
		rotationRepository   memory.RotationRepository
		statisticsRepository memory.StatisticsRepository
		slotID               int
		groupID              int
		err                  error
		expectedBannerID     int
	}{
		{
			rotationRepository:   memory.RotationRepository{},
			statisticsRepository: memory.StatisticsRepository{},
			slotID:               1,
			groupID:              1,
			err:                  ErrRotationsListEmpty,
			expectedBannerID:     0,
		},
		{
			rotationRepository: memory.RotationRepository{
				DB: map[int]repository.Rotation{
					1: {
						ID:          1,
						BannerID:    1,
						SlotID:      1,
						Description: "Banner 1",
						CreatedAt:   time.Now().UTC(),
					},
					2: {
						ID:          2,
						BannerID:    2,
						SlotID:      1,
						Description: "Banner 2",
						CreatedAt:   time.Now().UTC(),
					},
				},
				ID: 3,
			},
			statisticsRepository: memory.StatisticsRepository{
				DB: map[int]repository.Statistics{
					1: {
						ID:        1,
						Type:      repository.StatisticsTypeView,
						BannerID:  1,
						SlotID:    1,
						GroupID:   1,
						CreatedAt: time.Now().UTC(),
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
						CreatedAt:   time.Now().UTC(),
					},
					2: {
						ID:          2,
						BannerID:    2,
						SlotID:      1,
						Description: "Banner 2",
						CreatedAt:   time.Now().UTC(),
					},
					3: {
						ID:          3,
						BannerID:    3,
						SlotID:      1,
						Description: "Banner 3",
						CreatedAt:   time.Now().UTC(),
					},
				},
				ID: 4,
			},
			statisticsRepository: memory.StatisticsRepository{
				DB: map[int]repository.Statistics{
					1: {
						ID:        1,
						Type:      repository.StatisticsTypeView,
						BannerID:  1,
						SlotID:    1,
						GroupID:   1,
						CreatedAt: time.Now().UTC(),
					},
					2: {
						ID:        2,
						Type:      repository.StatisticsTypeClick,
						BannerID:  1,
						SlotID:    1,
						GroupID:   1,
						CreatedAt: time.Now().UTC(),
					},
					3: {
						ID:        3,
						Type:      repository.StatisticsTypeView,
						BannerID:  2,
						SlotID:    1,
						GroupID:   1,
						CreatedAt: time.Now().UTC(),
					},
					4: {
						ID:        4,
						Type:      repository.StatisticsTypeView,
						BannerID:  3,
						SlotID:    1,
						GroupID:   1,
						CreatedAt: time.Now().UTC(),
					},
					5: {
						ID:        5,
						Type:      repository.StatisticsTypeClick,
						BannerID:  3,
						SlotID:    1,
						GroupID:   1,
						CreatedAt: time.Now().UTC(),
					},
					6: {
						ID:        6,
						Type:      repository.StatisticsTypeView,
						BannerID:  1,
						SlotID:    1,
						GroupID:   1,
						CreatedAt: time.Now().UTC(),
					},
					7: {
						ID:        7,
						Type:      repository.StatisticsTypeView,
						BannerID:  3,
						SlotID:    1,
						GroupID:   1,
						CreatedAt: time.Now().UTC(),
					},
					8: {
						ID:        8,
						Type:      repository.StatisticsTypeClick,
						BannerID:  3,
						SlotID:    1,
						GroupID:   1,
						CreatedAt: time.Now().UTC(),
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
						CreatedAt:   time.Now().UTC(),
					},
					2: {
						ID:          2,
						BannerID:    2,
						SlotID:      2,
						Description: "Banner 2",
						CreatedAt:   time.Now().UTC(),
					},
					3: {
						ID:          3,
						BannerID:    3,
						SlotID:      3,
						Description: "Banner 3",
						CreatedAt:   time.Now().UTC(),
					},
					4: {
						ID:          4,
						BannerID:    4,
						SlotID:      2,
						Description: "Banner 4",
						CreatedAt:   time.Now().UTC(),
					},
				},
				ID: 5,
			},
			statisticsRepository: memory.StatisticsRepository{
				DB: map[int]repository.Statistics{
					1: {
						ID:        1,
						Type:      repository.StatisticsTypeView,
						BannerID:  1,
						SlotID:    1,
						GroupID:   1,
						CreatedAt: time.Now().UTC(),
					},
					2: {
						ID:        2,
						Type:      repository.StatisticsTypeClick,
						BannerID:  1,
						SlotID:    1,
						GroupID:   1,
						CreatedAt: time.Now().UTC(),
					},
					3: {
						ID:        3,
						Type:      repository.StatisticsTypeView,
						BannerID:  2,
						SlotID:    2,
						GroupID:   4,
						CreatedAt: time.Now().UTC(),
					},
					4: {
						ID:        4,
						Type:      repository.StatisticsTypeView,
						BannerID:  3,
						SlotID:    1,
						GroupID:   1,
						CreatedAt: time.Now().UTC(),
					},
					5: {
						ID:        5,
						Type:      repository.StatisticsTypeClick,
						BannerID:  3,
						SlotID:    1,
						GroupID:   1,
						CreatedAt: time.Now().UTC(),
					},
					6: {
						ID:        6,
						Type:      repository.StatisticsTypeView,
						BannerID:  1,
						SlotID:    1,
						GroupID:   1,
						CreatedAt: time.Now().UTC(),
					},
					7: {
						ID:        7,
						Type:      repository.StatisticsTypeView,
						BannerID:  4,
						SlotID:    2,
						GroupID:   4,
						CreatedAt: time.Now().UTC(),
					},
					8: {
						ID:        8,
						Type:      repository.StatisticsTypeClick,
						BannerID:  4,
						SlotID:    2,
						GroupID:   4,
						CreatedAt: time.Now().UTC(),
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
			StatisticsService: &StatisticsService{
				StatisticsRepository: &testCase.statisticsRepository,
			},
			StatisticsRepository: &testCase.statisticsRepository,
		}

		bannerID, _, err := rotationService.SelectBanner(context.Background(), testCase.slotID, testCase.groupID)

		if err != nil {
			assert.Error(t, testCase.err, &err)
		} else {
			assert.Equal(t, testCase.expectedBannerID, bannerID, "banners ids must match")
		}
	}
}
