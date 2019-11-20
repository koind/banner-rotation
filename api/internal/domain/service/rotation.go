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
	StatisticService    StatisticServiceInterface
	RotationRepository  repository.RotationRepositoryInterface
	StatisticRepository repository.StatisticRepositoryInterface
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
) (*repository.Statistic, error) {
	statistic, err := b.StatisticService.Save(ctx, rotation, groupID, repository.StatisticTypeClick)
	if err != nil {
		return nil, errors.Wrap(err, "error when set the transition")
	}

	return statistic, nil
}

// Selects a banner to display
func (b *RotationService) SelectBanner(
	ctx context.Context,
	slotID int,
	groupID int,
) (int, *repository.Statistic, error) {
	rotations, err := b.RotationRepository.FindAllBySlotID(ctx, slotID)
	if err != nil {
		return 0, nil, errors.Wrap(err, "error when searching for rotations by slot id for banner selection")
	}

	statistics, err := b.StatisticRepository.FindAllBySlotIDAndGroupID(ctx, slotID, groupID)
	if err != nil {
		return 0, nil, errors.Wrap(err, "error getting statistics for a selection of banner")
	}

	rotation, err := b.defineBanner(rotations, statistics)
	if err != nil {
		return 0, nil, errors.Wrap(err, "error while banner definition")
	}

	statistic, err := b.StatisticService.Save(ctx, *rotation, groupID, repository.StatisticTypeView)
	if err != nil {
		return 0, nil, errors.Wrap(err, "error while save view")
	}

	return rotation.BannerID, statistic, nil
}

// Determines which banner should be displayed
func (b *RotationService) defineBanner(
	rotations []*repository.Rotation,
	statistics []*repository.Statistic,
) (*repository.Rotation, error) {
	if len(rotations) <= 0 {
		return nil, ErrRotationsListEmpty
	}

	banners := make(map[int]repository.Banner, len(rotations))
	for _, rotation := range rotations {
		banners[rotation.BannerID] = repository.Banner{ID: rotation.BannerID}
	}

	for _, statistic := range statistics {
		banner, has := banners[statistic.BannerID]

		if !has {
			continue
		}

		if statistic.IsTypeView() {
			banner.Views++
		}

		if statistic.IsTypeClick() {
			banner.Clicks++
		}

		banner.GroupID = statistic.GroupID

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
