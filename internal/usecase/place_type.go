package usecase

import (
	"context"
	"errors"
	"fmt"

	"palback/internal/app"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
)

type PlaceTypeUseCase struct {
	repo app.PlaceTypeRepo
}

func NewPlaceTypeUseCase(repo app.PlaceTypeRepo) *PlaceTypeUseCase {
	return &PlaceTypeUseCase{
		repo: repo,
	}
}

func (c *PlaceTypeUseCase) Get(ctx context.Context, id int) (*model.PlaceType, error) {
	result, err := c.repo.Get(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, localErrors.ErrNotFound):
			return nil, app.ErrPlaceTypeNotFound
		default:
			return nil, fmt.Errorf("ошибка получения типа святого места по id: %w", err)
		}
	}

	return result, nil
}

func (c *PlaceTypeUseCase) GetAll(ctx context.Context) ([]model.PlaceType, error) {
	result, err := c.repo.GetAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка типов святых мест: %w", err)
	}

	return result, nil
}
