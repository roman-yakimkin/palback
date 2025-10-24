package port

import (
	"context"

	"palback/internal/domain/model"
)

type CountryRepo interface {
	Get(context.Context, string) (*model.Country, error)
	GetAll(context.Context) ([]model.Country, error)
	Create(context.Context, model.Country) (*model.Country, error)
	Update(context.Context, string, model.Country) error
	Delete(context.Context, string) error
	Order(context.Context, []string) error
}

type RegionRepo interface {
	Get(context.Context, int) (*model.Region, error)
	GetByCountry(context.Context, string) ([]model.Region, error)
	Create(context.Context, model.Region) (*model.Region, error)
	Update(context.Context, int, model.Region) error
	Delete(context.Context, int) error
}

type CityTypeRepo interface {
	Get(context.Context, int) (*model.CityType, error)
	GetAll(context.Context) ([]model.CityType, error)
}

type PlaceTypeRepo interface {
	Get(context.Context, int) (*model.PlaceType, error)
	GetAll(context.Context) ([]model.PlaceType, error)
}
