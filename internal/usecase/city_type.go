package usecase

import (
	"context"
	"errors"
	"fmt"

	"palback/internal/app"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
)

type CityTypeUseCase struct {
	repo app.CityTypeRepo
}

func NewCityTypeUseCase(repo app.CityTypeRepo) *CityTypeUseCase {
	return &CityTypeUseCase{
		repo: repo,
	}
}

func (c *CityTypeUseCase) Get(ctx context.Context, id int) (*model.CityType, error) {
	result, err := c.repo.Get(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, localErrors.ErrNotFound):
			return nil, app.ErrCityTypeNotFound
		default:
			return nil, fmt.Errorf("ошибка получения типа населенного пункта по id: %w", err)
		}
	}

	return result, nil
}

func (c *CityTypeUseCase) GetAll(ctx context.Context) ([]model.CityType, error) {
	result, err := c.repo.GetAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка типов населенных пунктов: %w", err)
	}

	return result, nil
}
