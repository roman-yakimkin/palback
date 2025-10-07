package domain

import (
	"context"
	"errors"
	"palback/internal/domain/model"
)

var (
	ErrCountryNotFound      = errors.New("страна не найдена")
	ErrCountryAlreadyAdded  = errors.New("страна с таким id уже добавлена")
	ErrCountryNameNotUnique = errors.New("страна должна иметь уникальное название")
)

type CountryService interface {
	Get(ctx context.Context, id string) (*model.Country, error)
	GetAll(ctx context.Context) ([]model.Country, error)
	Create(ctx context.Context, country model.Country) (*model.Country, error)
	Update(ctx context.Context, id string, country model.Country) error
	Delete(ctx context.Context, id string) error
}

var (
	ErrCountryHasNotRegions = errors.New("страна не имеет регионов")
	ErrRegionNotFound       = errors.New("регион не найден")
	ErrRegionNotUnique      = errors.New("в пределах одной страны регион должен иметь уникальное название")
)

type RegionService interface {
	Get(ctx context.Context, id int) (*model.Region, error)
	GetByCountry(ctx context.Context, countryId string) ([]model.Region, error)
	Create(ctx context.Context, region model.Region) (*model.Region, error)
	Update(ctx context.Context, id int, region model.Region) error
	Delete(ctx context.Context, id int) error
}

var (
	ErrCityTypeNotFound = errors.New("тип населенного пункта не найден")
)

type CityTypeService interface {
	Get(ctx context.Context, id int) (*model.CityType, error)
	GetAll(ctx context.Context) ([]model.CityType, error)
}
