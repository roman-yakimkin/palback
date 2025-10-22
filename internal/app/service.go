package app

import (
	"context"
	"errors"

	appModel "palback/internal/app/model"
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
	Order(ctx context.Context, ids []string) error
}

var (
	ErrCountryHasNotRegions = errors.New("страна не имеет регионов")
	ErrRegionNotFound       = errors.New("регион не найден")
	ErrRegionNotUnique      = errors.New("в пределах одной страны регион должен иметь уникальное название")
)

type RegionService interface {
	Get(ctx context.Context, id int) (*appModel.RegionDetail, error)
	GetByCountry(ctx context.Context, countryId string) (appModel.RegionList, error)
	Create(ctx context.Context, region model.Region) (*appModel.RegionDetail, error)
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

var (
	ErrPlaceTypeNotFound = errors.New("тип святого места не найден")
)

type PlaceTypeService interface {
	Get(ctx context.Context, id int) (*model.PlaceType, error)
	GetAll(ctx context.Context) ([]model.PlaceType, error)
}
