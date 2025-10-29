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

type UserRepo interface {
	Get(context.Context, int) (*model.User, error)
	GetByIdentifier(ctx context.Context, identifier string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetAll(context.Context) ([]model.User, error)
	Create(context.Context, model.User) (*model.User, error)
	Delete(context.Context, int) error
	UpdateEmailVerified(ctx context.Context, email string) error
	UpdatePassword(ctx context.Context, email, hashedPassword string) error
	IncrementSessionVersion(ctx context.Context, email string) error
}

type RoleRepo interface {
	Get(context.Context, model.RoleID) (*model.Role, error)
	GetAll(context.Context) ([]model.Role, error)
	GetAllMap(context.Context) (map[model.RoleID]*model.Role, error)
}
