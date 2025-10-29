package dto

import "palback/internal/domain/model"

type RoleResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func CreateRoleResponse(src model.Role) RoleResponse {
	return RoleResponse{
		ID:   string(src.ID),
		Name: src.Name,
	}
}
