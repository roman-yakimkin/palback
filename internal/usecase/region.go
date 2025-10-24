package usecase

import (
	"context"
	"errors"
	"fmt"

	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
	"palback/internal/pkg/helpers"
	ucModel "palback/internal/usecase/model"
	"palback/internal/usecase/port"
)

type RegionUseCase struct {
	countryService CountryService
	repo           port.RegionRepo
}

func NewRegionUseCase(countryService CountryService, repo port.RegionRepo) *RegionUseCase {
	return &RegionUseCase{
		countryService: countryService,
		repo:           repo,
	}
}

func (s *RegionUseCase) Get(ctx context.Context, id int) (*ucModel.RegionDetail, error) {
	region, err := s.repo.Get(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("ошибка получения региона по id: %w", err)
	}

	country, err := s.countryService.Get(ctx, region.CountryID)

	if err != nil {
		return nil, fmt.Errorf("ошибка получения страны по id: %w", err)
	}

	regionDetail := ucModel.CreateRegionDetail(helpers.FromPtr(region), helpers.FromPtr(country))

	return &regionDetail, nil
}

func (s *RegionUseCase) GetByCountry(ctx context.Context, countryID string) (result ucModel.RegionList, err error) {
	country, err := s.countryService.Get(ctx, countryID)
	if err != nil {
		switch {
		case errors.Is(err, ErrCountryNotFound):
			return result, fmt.Errorf("ошибка получения страны: %w", err)
		default:
			return result, fmt.Errorf("ошибка проверки страны на существование: %w", err)
		}
	}

	regions, err := s.repo.GetByCountry(ctx, countryID)

	if err != nil {
		return result, fmt.Errorf("ошибка при получении списка регионов: %w", err)
	}

	result = ucModel.CreateRegionList(regions, []model.Country{helpers.FromPtr(country)})

	return result, nil
}

func (s *RegionUseCase) Create(ctx context.Context, region model.Region) (*ucModel.RegionDetail, error) {
	country, err := s.countryService.Get(ctx, region.CountryID)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки страны на возможность добавления регионов: %w", err)
	}

	if !country.HasRegions {
		return nil, ErrCountryHasNotRegions
	}

	reg, err := s.repo.Create(ctx, region)

	if err != nil {
		switch {
		default:
			return nil, fmt.Errorf("ошибка добавления региона: %w", err)
		}
	}

	result := ucModel.CreateRegionDetail(helpers.FromPtr(reg), helpers.FromPtr(country))

	return &result, nil
}

func (s *RegionUseCase) Update(ctx context.Context, id int, region model.Region) error {
	country, err := s.countryService.Get(ctx, region.CountryID)
	if err != nil {
		return fmt.Errorf("ошибка проверки страны на возможность добавления регионов: %w", err)
	}

	if !country.HasRegions {
		return ErrCountryHasNotRegions
	}

	err = s.repo.Update(ctx, id, region)

	if err != nil {
		switch {
		case errors.Is(err, localErrors.ErrNotFound):
			return ErrRegionNotFound
		default:
			return fmt.Errorf("ошибка обновления региона: %w", err)
		}
	}

	return nil
}

func (s *RegionUseCase) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)

	if errors.Is(err, localErrors.ErrNotFound) {
		return ErrRegionNotFound
	}

	if err != nil {
		return fmt.Errorf("ошибка удаления региона: %w", err)
	}

	return nil
}
