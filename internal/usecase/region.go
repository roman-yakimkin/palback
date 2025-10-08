package usecase

import (
	"context"
	"errors"
	"fmt"
	
	"palback/internal/domain"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
)

type RegionUseCase struct {
	countryService domain.CountryService
	repo           domain.RegionRepo
}

func NewRegionUseCase(countryService domain.CountryService, repo domain.RegionRepo) *RegionUseCase {
	return &RegionUseCase{
		countryService: countryService,
		repo:           repo,
	}
}

func (s *RegionUseCase) Get(ctx context.Context, id int) (*model.Region, error) {
	region, err := s.repo.Get(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("ошибка получения региона по id: %w", err)
	}

	return region, nil
}

func (s *RegionUseCase) GetByCountry(ctx context.Context, countryID string) ([]model.Region, error) {
	_, err := s.countryService.Get(ctx, countryID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCountryNotFound):
			return nil, fmt.Errorf("ошибка получения страны: %w", err)
		default:
			return nil, fmt.Errorf("ошибка проверки страны на существование: %w", err)
		}
	}

	regions, err := s.repo.GetByCountry(ctx, countryID)

	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка регионов: %w", err)
	}

	return regions, nil
}

func (s *RegionUseCase) Create(ctx context.Context, region model.Region) (*model.Region, error) {
	country, err := s.countryService.Get(ctx, region.CountryID)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки страны на возможность добавления регионов: %w", err)
	}

	if !country.HasRegions {
		return nil, domain.ErrCountryHasNotRegions
	}

	result, err := s.repo.Create(ctx, region)

	if err != nil {
		switch {
		default:
			return nil, fmt.Errorf("ошибка добавления региона: %w", err)
		}
	}

	return result, nil
}

func (s *RegionUseCase) Update(ctx context.Context, id int, region model.Region) error {
	country, err := s.countryService.Get(ctx, region.CountryID)
	if err != nil {
		return fmt.Errorf("ошибка проверки страны на возможность добавления регионов: %w", err)
	}

	if !country.HasRegions {
		return domain.ErrCountryHasNotRegions
	}

	err = s.repo.Update(ctx, id, region)

	if err != nil {
		switch {
		case errors.Is(err, localErrors.ErrNotFound):
			return domain.ErrRegionNotFound
		default:
			return fmt.Errorf("ошибка обновления региона: %w", err)
		}
	}

	return nil
}

func (s *RegionUseCase) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)

	if errors.Is(err, localErrors.ErrNotFound) {
		return domain.ErrRegionNotFound
	}

	if err != nil {
		return fmt.Errorf("ошибка удаления региона: %w", err)
	}

	return nil
}
