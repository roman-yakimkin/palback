package http

import (
	"context"
	"net/http"
	"palback/internal/domain/model"
)

// Authenticator управляет сессиями пользователя.
// Реализация зависит от транспорта (HTTP cookies, JWT и т.д.).
type Authenticator interface {
	// Login создаёт сессию для пользователя и устанавливает cookie.
	Login(ctx context.Context, w http.ResponseWriter, r *http.Request, user *model.User) error

	// GetUserID извлекает ID пользователя из сессии.
	// Возвращает ошибку, если сессия отсутствует или недействительна.
	GetUserID(ctx context.Context, r *http.Request) (int, error)

	// Logout удаляет сессию и cookie.
	Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}
