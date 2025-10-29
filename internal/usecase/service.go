package usecase

import (
	"context"
	"palback/internal/domain/model"
	ucModel "palback/internal/usecase/model"
)

type CountryService interface {
	Get(ctx context.Context, id string) (*model.Country, error)
	GetAll(ctx context.Context) ([]model.Country, error)
	Create(ctx context.Context, country model.Country) (*model.Country, error)
	Update(ctx context.Context, id string, country model.Country) error
	Delete(ctx context.Context, id string) error
	Order(ctx context.Context, ids []string) error
}

type RegionService interface {
	Get(ctx context.Context, id int) (*ucModel.RegionDetail, error)
	GetByCountry(ctx context.Context, countryId string) (ucModel.RegionList, error)
	Create(ctx context.Context, region model.Region) (*ucModel.RegionDetail, error)
	Update(ctx context.Context, id int, region model.Region) error
	Delete(ctx context.Context, id int) error
}

type CityTypeService interface {
	Get(ctx context.Context, id int) (*model.CityType, error)
	GetAll(ctx context.Context) ([]model.CityType, error)
}

type PlaceTypeService interface {
	Get(ctx context.Context, id int) (*model.PlaceType, error)
	GetAll(ctx context.Context) ([]model.PlaceType, error)
}

type RoleService interface {
	Get(ctx context.Context, id model.RoleID) (*model.Role, error)
	GetAll(ctx context.Context) ([]model.Role, error)
	GetAllMap(ctx context.Context) (map[model.RoleID]*model.Role, error)
}

type UserService interface {
	Get(ctx context.Context, id int) (*ucModel.UserDetail, error)
	GetAll(ctx context.Context) (ucModel.UserList, error)
	Register(ctx context.Context, userName, email, password string) (*ucModel.UserDetail, error)
	VerifyEmail(ctx context.Context, token string) error
	ResendVerificationEmail(ctx context.Context, email string) error
	Login(ctx context.Context, identifier, password string) (*ucModel.UserDetail, error)
	RequestPasswordReset(ctx context.Context, email string) error
	ConfirmPasswordReset(ctx context.Context, token, newPassword string) error
}
