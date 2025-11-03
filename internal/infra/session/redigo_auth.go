package session

import (
	"context"
	"net/http"
	"palback/internal/config"
	"palback/internal/domain/model"
	"palback/internal/usecase"

	"github.com/boj/redistore"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/sessions"
)

type UserGetterRepo interface {
	Get(ctx context.Context, id int) (*model.User, error)
}

type RedigoAuthenticator struct {
	store sessions.Store
	repo  UserGetterRepo
}

func NewRedigoAuthenticator(
	cfg *config.Config,
	redisPool *redis.Pool,
	repo UserGetterRepo,
) (*RedigoAuthenticator, error) {
	store, err := redistore.NewRediStoreWithPool(redisPool, []byte(cfg.RedisSecretKey))
	if err != nil {
		return nil, err
	}

	// Настройка cookie
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * cfg.SessionDays, // 24 часа
		HttpOnly: true,
		Secure:   cfg.IsProduction,
		SameSite: http.SameSiteLaxMode,
	}

	return &RedigoAuthenticator{
		store: store,
		repo:  repo,
	}, nil
}

func (a *RedigoAuthenticator) Login(_ context.Context, w http.ResponseWriter, r *http.Request, user *model.User) error {
	session, err := a.store.Get(r, "session")
	if err != nil {
		return err
	}
	session.Values["user_id"] = user.ID
	session.Values["session_version"] = user.SessionVersion
	return session.Save(r, w)
}

func (a *RedigoAuthenticator) GetUserID(ctx context.Context, r *http.Request) (int, error) {
	session, err := a.store.Get(r, "session")

	if err != nil {
		return 0, err
	}

	userID, ok1 := session.Values["user_id"].(int)
	storedSession, ok2 := session.Values["session_version"]

	if !ok1 || !ok2 || userID <= 0 {
		return 0, usecase.ErrUnauthenticated
	}

	userInfo, err := a.repo.Get(ctx, userID)
	if err != nil {
		return 0, usecase.ErrUnauthenticated
	}

	if userInfo.SessionVersion != storedSession {
		return 0, usecase.ErrSessionExpired
	}

	return userID, nil
}

func (a *RedigoAuthenticator) Logout(_ context.Context, w http.ResponseWriter, r *http.Request) error {
	session, err := a.store.Get(r, "session")
	if err != nil {
		return err
	}

	// Удаляем cookie
	session.Options.MaxAge = -1
	session.Values = make(map[any]any)
	return session.Save(r, w)
}
