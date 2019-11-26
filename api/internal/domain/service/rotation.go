package service

import (
	"context"
	"github.com/koind/banner-rotation/api/internal/algorithm"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"github.com/pkg/errors"
)

var (
	ErrRotationsListEmpty = errors.New("rotations list can't be empty")
)

// Rotation service
type RotationService struct {
	StatisticsService    StatisticsServiceInterface
	RotationRepository   repository.RotationRepositoryInterface
	StatisticsRepository repository.StatisticsRepositoryInterface
}

// Adds a new banner to the rotation
func (b *RotationService) Add(ctx context.Context, rotation repository.Rotation) (*repository.Rotation, error) {
	newRotation, err := b.RotationRepository.Add(ctx, rotation)
	if err != nil {
		return nil, errors.Wrap(err, "error when adding banner in the rotation")
	}

	return newRotation, nil
}

// Removes the banner from the rotation
func (b *RotationService) Remove(ctx context.Context, bannerID int) error {
	err := b.RotationRepository.Remove(ctx, bannerID)
	if err != nil {
		return errors.Wrap(err, "error while removing banner from rotation")
	}

	return nil
}

// Increases the jump count by 1 for the specified banner in the specified group
func (b *RotationService) SetTransition(
	ctx context.Context,
	rotation repository.Rotation,
	groupID int,
) (*repository.Statistics, error) {
	statistics, err := b.StatisticsService.Save(ctx, rotation, groupID, repository.StatisticsTypeClick)
	if err != nil {
		return nil, errors.Wrap(err, "error when set the transition")
	}

	return statistics, nil
}

// Selects a banner to display
func (b *RotationService) SelectBanner(
	ctx context.Context,
	slotID int,
	groupID int,
) (int, *repository.Statistics, error) {
	rotations, err := b.RotationRepository.FindAllBySlotID(ctx, slotID)
	if err != nil {
		return 0, nil, errors.Wrap(err, "error when searching for rotations by slot id for banner selection")
	}

	statisticsList, err := b.StatisticsRepository.FindAllBySlotIDAndGroupID(ctx, slotID, groupID)
	if err != nil {
		return 0, nil, errors.Wrap(err, "error getting statistics for a selection of banner")
	}

	rotation, err := b.defineBanner(rotations, statisticsList)
	if err != nil {
		return 0, nil, errors.Wrap(err, "error while banner definition")
	}

	statistics, err := b.StatisticsService.Save(ctx, *rotation, groupID, repository.StatisticsTypeView)
	if err != nil {
		return 0, nil, errors.Wrap(err, "error while save view")
	}

	return rotation.BannerID, statistics, nil
}

// Determines which banner should be displayed
func (b *RotationService) defineBanner(
	rotations []*repository.Rotation,
	statisticsList []*repository.Statistics,
) (*repository.Rotation, error) {
	if len(rotations) <= 0 {
		return nil, ErrRotationsListEmpty
	}

	banners := make(map[int]repository.Banner, len(rotations))
	for _, rotation := range rotations {
		banners[rotation.BannerID] = repository.Banner{ID: rotation.BannerID}
	}

	for _, statistics := range statisticsList {
		banner, has := banners[statistics.BannerID]

		if !has {
			continue
		}

		if statistics.IsTypeView() {
			banner.Views++
		}

		if statistics.IsTypeClick() {
			banner.Clicks++
		}

		banner.GroupID = statistics.GroupID

		banners[banner.ID] = banner
	}

	selected := make([]int, 0, len(banners))
	reward := make([]float64, 0, len(banners))
	arms := make(map[int]int, len(banners))
	i := 0

	for bannerID, banner := range banners {
		arms[i] = bannerID
		selected = append(selected, banner.Views)
		reward = append(reward, banner.Clicks)
		i++
	}

	ucb1, err := algorithm.NewUCB1(selected, reward)
	if err != nil {
		return nil, err
	}

	arm := ucb1.SelectArm()
	bannerID := arms[arm]
	rotation := new(repository.Rotation)

	for _, r := range rotations {
		if bannerID == r.BannerID {
			rotation = r
		}
	}

	return rotation, nil
}
