package usecase

import (
	"context"
	"errors"
	"fmt"

	"palback/internal/app"
	appModel "palback/internal/app/model"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
	"palback/internal/pkg/helpers"
)

type RegionUseCase struct {
	countryService app.CountryService
	repo           app.RegionRepo
}

func NewRegionUseCase(countryService app.CountryService, repo app.RegionRepo) *RegionUseCase {
	return &RegionUseCase{
		countryService: countryService,
		repo:           repo,
	}
}

func (s *RegionUseCase) Get(ctx context.Context, id int) (*appModel.RegionDetail, error) {
	region, err := s.repo.Get(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("ошибка получения региона по id: %w", err)
	}

	country, err := s.countryService.Get(ctx, region.CountryID)

	if err != nil {
		return nil, fmt.Errorf("ошибка получения страны по id: %w", err)
	}

	regionDetail := appModel.CreateRegionDetail(helpers.FromPtr(region), helpers.FromPtr(country))

	return &regionDetail, nil
}

func (s *RegionUseCase) GetByCountry(ctx context.Context, countryID string) (result appModel.RegionList, err error) {
	country, err := s.countryService.Get(ctx, countryID)
	if err != nil {
		switch {
		case errors.Is(err, app.ErrCountryNotFound):
			return result, fmt.Errorf("ошибка получения страны: %w", err)
		default:
			return result, fmt.Errorf("ошибка проверки страны на существование: %w", err)
		}
	}

	regions, err := s.repo.GetByCountry(ctx, countryID)

	if err != nil {
		return result, fmt.Errorf("ошибка при получении списка регионов: %w", err)
	}

	result = appModel.CreateRegionList(regions, []model.Country{helpers.FromPtr(country)})

	return result, nil
}

func (s *RegionUseCase) Create(ctx context.Context, region model.Region) (*appModel.RegionDetail, error) {
	country, err := s.countryService.Get(ctx, region.CountryID)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки страны на возможность добавления регионов: %w", err)
	}

	if !country.HasRegions {
		return nil, app.ErrCountryHasNotRegions
	}

	reg, err := s.repo.Create(ctx, region)

	if err != nil {
		switch {
		default:
			return nil, fmt.Errorf("ошибка добавления региона: %w", err)
		}
	}

	result := appModel.CreateRegionDetail(helpers.FromPtr(reg), helpers.FromPtr(country))

	return &result, nil
}

func (s *RegionUseCase) Update(ctx context.Context, id int, region model.Region) error {
	country, err := s.countryService.Get(ctx, region.CountryID)
	if err != nil {
		return fmt.Errorf("ошибка проверки страны на возможность добавления регионов: %w", err)
	}

	if !country.HasRegions {
		return app.ErrCountryHasNotRegions
	}

	err = s.repo.Update(ctx, id, region)

	if err != nil {
		switch {
		case errors.Is(err, localErrors.ErrNotFound):
			return app.ErrRegionNotFound
		default:
			return fmt.Errorf("ошибка обновления региона: %w", err)
		}
	}

	return nil
}

func (s *RegionUseCase) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)

	if errors.Is(err, localErrors.ErrNotFound) {
		return app.ErrRegionNotFound
	}

	if err != nil {
		return fmt.Errorf("ошибка удаления региона: %w", err)
	}

	return nil
}
