package repository

import (
	"context"
	"palback/internal/domain/model"
)

type RoleRepo struct {
	roles    []model.Role
	rolesMap map[model.RoleID]*model.Role
}

func NewRoleRepo() *RoleRepo {
	roles := []model.Role{
		{
			ID:   model.RoleAdmin,
			Name: "администратор",
		},
		{
			ID:   model.RoleUser,
			Name: "пользователь",
		},
	}

	rolesMap := make(map[model.RoleID]*model.Role)
	rolesMap[model.RoleAdmin] = &roles[0]
	rolesMap[model.RoleUser] = &roles[1]

	return &RoleRepo{
		roles:    roles,
		rolesMap: rolesMap,
	}
}

func (r *RoleRepo) Get(_ context.Context, id model.RoleID) (*model.Role, error) {
	return r.rolesMap[id], nil
}

func (r *RoleRepo) GetAll(_ context.Context) ([]model.Role, error) {
	return r.roles, nil
}

func (r *RoleRepo) GetAllMap(_ context.Context) (map[model.RoleID]*model.Role, error) {
	return r.rolesMap, nil
}
