package repository

import "context"

// Модель баннера
type Banner struct {
	ID          int    `json:"id" db:"id"`
	SlotID      int    `json:"slot_id" db:"slot_id"`
	Description string `json:"description" db:"description"`
}

//
type BannerRepositoryInterface interface {
	// Добавляет новый баннер в ротацию в данном слоте
	Add(ctx context.Context, banner Banner) (*Banner, error)

	// Удаляет баннер из ротации
	Remove(ctx context.Context, ID int) error

	// Увеличивает счетчик переходов на 1 для указанного баннера в указанной группе
	IncreaseCountTransition(ctx context.Context, bannerID int, GroupID int) error
}
