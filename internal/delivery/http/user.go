package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"palback/internal/delivery/http/dto"
	"palback/internal/domain/model"
	"palback/internal/pkg/helpers"
	"palback/internal/usecase"
	"palback/internal/usecase/port"
)

type UserHandler struct {
	service     usecase.UserService
	auth        Authenticator
	rateLimiter port.RateLimiter
}

func NewUserHandler(service usecase.UserService, auth Authenticator, rateLimiter port.RateLimiter) *UserHandler {
	return &UserHandler{
		service:     service,
		auth:        auth,
		rateLimiter: rateLimiter,
	}
}

func (h *UserHandler) Register(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	data, err := h.service.Register(ctx, req.Username, email, req.Password)
	if err != nil {
		switch {
		default:
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf("ошибка регистрации пользователя: %s", err.Error()),
			)
		}
	}

	dataRec := helpers.FromPtr(data)

	c.Response().Header().Set("location", "/users/"+strconv.Itoa(dataRec.ID))

	return c.JSON(http.StatusCreated, dto.CreateUserResponse(dataRec))
}

func (h *UserHandler) VerifyEmail(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.VerifyEmailRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "неверный json")
	}

	if req.Token == "" {
		return c.JSON(http.StatusBadRequest, "отсутсвует токен")
	}

	err := h.service.VerifyEmail(ctx, req.Token)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNoReplyFromKeyValueStorage):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Errorf("ошибка при проверке e-mail: %w", err),
			)
		}
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "e-mail успешно подтвержден"})
}

func (h *UserHandler) ResendVerification(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.ResendVerificationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "неверный json")
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	if email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "e-mail не задан")
	}

	allowed, err := h.rateLimiter.Allow(c.Request().Context(), "rl:resend-verification:"+email, 5, 3600)
	if err != nil {
		log.Warn("rate limit check failed", "email", email, "error", err)
	}

	if !allowed {
		return echo.NewHTTPError(
			http.StatusTooManyRequests,
			"слишком много запросов на изменение пароля от данного e-mail, подождите некоторое время",
		)
	}

	err = h.service.ResendVerificationEmail(ctx, email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "внутренная ошибка сервиса")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "если вы зарегистрировались на сайте, но пока не подтвердили e-mail, новое письмо для подтверждения отправлено на вашу почту",
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "неверный json")
	}

	if req.Identifier == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "логин или пароль не заданы")
	}

	data, err := h.service.Login(ctx, req.Identifier, req.Password)
	if err != nil || data == nil {
		switch {
		case errors.Is(err, usecase.ErrUncheckedEmail):
			return echo.NewHTTPError(http.StatusForbidden, err)
		case errors.Is(err, usecase.ErrUserInvalidCredentials):
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		default:
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf("невозможно войти в систему: %s", err),
			)
		}
	}

	err = h.auth.Login(ctx, c.Response(), c.Request(), &model.User{
		ID:             data.ID,
		SessionVersion: data.SessionVersion,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка сессии")
	}

	dataRec := helpers.FromPtr(data)

	return c.JSON(http.StatusOK, dto.CreateUserResponse(dataRec))
}

func (h *UserHandler) Logout(c echo.Context) error {
	ctx := c.Request().Context()

	err := h.auth.Logout(ctx, c.Response(), c.Request())
	if err != nil {
		log.Error("logout failed", err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "успешный выход",
	})
}

func (h *UserHandler) ResetPassword(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.ResetPasswordRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "неверный json")
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	if email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "e-mail не задан")
	}

	// Rate limit по email
	allowed, err := h.rateLimiter.Allow(c.Request().Context(), "rl:reset:"+email, 5, 3600)
	if err != nil {
		log.Warn("rate limit check failed", "email", email, "error", err)
	}

	if !allowed {
		return echo.NewHTTPError(
			http.StatusTooManyRequests,
			"слишком много запросов на изменение пароля от данного e-mail, подождите некоторое время",
		)
	}

	err = h.service.RequestPasswordReset(ctx, email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка выполнения запроса")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "если ваш e-mail зарегистрирован, вы получите письмо со ссылкой на сброс пароля",
	})
}

func (h *UserHandler) ResetPasswordConfirm(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.ResetPasswordConfirmRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "неверный json")
	}

	if req.Token == "" || len(req.NewPassword) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "требуется токен и пароль длиной не менее 6 символов")
	}

	err := h.service.ConfirmPasswordReset(ctx, req.Token, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidToken):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "не удалось обновить пароль")
		}
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "пароль был успешно обновлён",
	})
}

func (h *UserHandler) Me(c echo.Context) error {
	ctx := c.Request().Context()

	userId, err := h.auth.GetUserID(ctx, c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "неверная сессия")
	}

	if userId <= 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "пользователь не авторизован")
	}

	user, err := h.service.Get(ctx, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "пользователь не найден")
	}

	return c.JSON(http.StatusOK, dto.CreateUserResponse(helpers.FromPtr(user)))
}

func (h *UserHandler) Profile(c echo.Context) error {
	return nil
}

func (h *UserHandler) Delete(c echo.Context) error {
	return nil
}
