package usecase

import (
	"context"
	"palback/internal/domain/model"
	"palback/internal/usecase/port"
)

type RoleUseCase struct {
	repo port.RoleRepo
}

func NewRoleUseCase(repo port.RoleRepo) *RoleUseCase {
	return &RoleUseCase{
		repo: repo,
	}
}

func (s *RoleUseCase) Get(ctx context.Context, id model.RoleID) (*model.Role, error) {
	return s.repo.Get(ctx, id)
}

func (s *RoleUseCase) GetAll(ctx context.Context) ([]model.Role, error) {
	return s.repo.GetAll(ctx)
}

func (s *RoleUseCase) GetAllMap(ctx context.Context) (map[model.RoleID]*model.Role, error) {
	return s.repo.GetAllMap(ctx)
}
