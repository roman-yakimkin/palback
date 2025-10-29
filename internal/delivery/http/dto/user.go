package dto

import (
	"time"

	ucModel "palback/internal/usecase/model"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"` // username или email
	Password   string `json:"password" validate:"required"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordConfirmRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type ResendVerificationRequest struct {
	Email string `json:"email"`
}

type UserResponse struct {
	ID        int          `json:"id"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	CreatedAt time.Time    `json:"created_at"`
	Role      RoleResponse `json:"role"`
}

func CreateUserResponse(src ucModel.UserDetail) UserResponse {
	return UserResponse{
		ID:        src.ID,
		Username:  src.Username,
		Email:     src.Email,
		CreatedAt: src.CreatedAt,
		Role:      CreateRoleResponse(src.Role),
	}
}
