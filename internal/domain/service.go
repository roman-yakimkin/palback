package domain

import (
	"context"
	"errors"
	"palback/internal/domain/model"
)

var (
	ErrCountryNotFound     = errors.New("страна не найдена")
	ErrCountryAlreadyAdded = errors.New("страна с таким id уже добавлена")
)

type CountryService interface {
	Get(ctx context.Context, id string) (*model.Country, error)
	GetAll(ctx context.Context) ([]model.Country, error)
	Post(ctx context.Context, country model.Country) (*model.Country, error)
	Put(ctx context.Context, id string, country model.Country) error
	Delete(ctx context.Context, id string) error
}
