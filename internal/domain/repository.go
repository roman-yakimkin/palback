package domain

import (
	"context"
	"palback/internal/domain/model"
)

type CountryRepo interface {
	Get(context.Context, string) (*model.Country, error)
	GetAll(context.Context) ([]model.Country, error)
	Post(context.Context, model.Country) (*model.Country, error)
	Put(context.Context, string, model.Country) error
	Delete(context.Context, string) error
}
