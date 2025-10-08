package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"palback/internal/domain"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
)

type CountryUseCase struct {
	repo domain.CountryRepo
}

func NewCountryUseCase(repo domain.CountryRepo) *CountryUseCase {
	return &CountryUseCase{
		repo: repo,
	}
}

func (c *CountryUseCase) Get(ctx context.Context, id string) (*model.Country, error) {
	country, err := c.repo.Get(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, localErrors.ErrNotFound):
			return nil, domain.ErrCountryNotFound
		default:
			return nil, fmt.Errorf("ошибка получения страны по id: %w", err)
		}
	}

	return country, nil
}

func (c *CountryUseCase) GetAll(ctx context.Context) ([]model.Country, error) {
	countries, err := c.repo.GetAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка стран: %w", err)
	}

	return countries, nil
}

func (c *CountryUseCase) Create(ctx context.Context, country model.Country) (*model.Country, error) {
	result, err := c.repo.Create(ctx, country)

	if err != nil {
		switch {
		default:
			return nil, fmt.Errorf("ошибка добавления страны: %w", err)
		}
	}

	return result, nil
}

func (c *CountryUseCase) Update(ctx context.Context, id string, country model.Country) error {
	err := c.repo.Update(ctx, id, country)

	if err != nil {
		switch {
		case errors.Is(err, localErrors.ErrNotFound):
			return domain.ErrCountryNotFound
		case strings.Contains(err.Error(), "duplicate key"):
			return domain.ErrCountryAlreadyAdded
		default:
			return fmt.Errorf("ошибка обновления страны: %w", err)
		}
	}

	return nil
}

func (c *CountryUseCase) Delete(ctx context.Context, id string) error {
	err := c.repo.Delete(ctx, id)

	if errors.Is(err, localErrors.ErrNotFound) {
		return domain.ErrCountryNotFound
	}

	if err != nil {
		return fmt.Errorf("ошибка удаления страны: %w", err)
	}

	return nil
}
