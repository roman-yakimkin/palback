package repository

import (
	"context"
	"database/sql"
	"errors"
	
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
)

type PlaceTypeRepo struct {
	db *sql.DB
}

func NewPlaceTypeRepo(db *sql.DB) *PlaceTypeRepo {
	return &PlaceTypeRepo{
		db: db,
	}
}

type placeTypeDTO struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Weight int    `json:"weight"`
}

func (dto *placeTypeDTO) ToModel() model.PlaceType {
	return model.PlaceType{
		ID:     dto.ID,
		Name:   dto.Name,
		Weight: dto.Weight,
	}
}

// Get Получить информацию об одном населенном пункте
func (r *PlaceTypeRepo) Get(ctx context.Context, id int) (*model.PlaceType, error) {
	q := `select id, name, weight from place_types where id = $1`

	var dto placeTypeDTO

	err := r.db.QueryRowContext(ctx, q, id).Scan(&dto.ID, &dto.Name, &dto.Weight)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, localErrors.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	placeType := dto.ToModel()
	return &placeType, nil
}

// GetAll Получить информацию обо всех населенных пунктах
func (r *PlaceTypeRepo) GetAll(ctx context.Context) ([]model.PlaceType, error) {
	q := `select id, name, weight from place_types order by weight`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var placeTypes []model.PlaceType

	for rows.Next() {
		var dto placeTypeDTO

		err := rows.Scan(&dto.ID, &dto.Name, &dto.Weight)
		if err != nil {
			return nil, err
		}

		placeTypes = append(placeTypes, dto.ToModel())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return placeTypes, nil
}
