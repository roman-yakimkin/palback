package repository

import (
	"context"
	"database/sql"
	"errors"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
	"palback/internal/usecase"
	"strings"
	"time"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

type userDTO struct {
	ID             int       `json:"id"`
	RoleID         string    `json:"role_id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	CreatedAt      time.Time `json:"created_at"`
	EmailVerified  bool      `json:"email_verified"`
	SessionVersion int64     `json:"session_version"`
}

func (dto *userDTO) ToModel() model.User {
	return model.User{
		ID:             dto.ID,
		RoleID:         model.RoleID(dto.RoleID),
		Username:       dto.Username,
		Email:          dto.Email,
		Password:       dto.Password,
		CreatedAt:      dto.CreatedAt,
		EmailVerified:  dto.EmailVerified,
		SessionVersion: dto.SessionVersion,
	}
}

func (r *UserRepo) Get(ctx context.Context, id int) (*model.User, error) {
	q := `
select id, role_id, username, email, password, created_at, email_verified, session_version 
from users 
where id = $1
`

	var dto userDTO

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&dto.ID,
		&dto.RoleID,
		&dto.Username,
		&dto.Email,
		&dto.Password,
		&dto.CreatedAt,
		&dto.EmailVerified,
		&dto.SessionVersion,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, localErrors.ErrNotFound
	}

	user := dto.ToModel()

	return &user, nil
}

func (r *UserRepo) GetByIdentifier(ctx context.Context, identifier string) (*model.User, error) {
	q := `
SELECT id, role_id, username, email, password, created_at, email_verified, session_version 
FROM users 
WHERE username = $1 OR email = $2
`

	var dto userDTO

	err := r.db.QueryRowContext(ctx, q, identifier, identifier).Scan(
		&dto.ID,
		&dto.RoleID,
		&dto.Username,
		&dto.Email,
		&dto.Password,
		&dto.CreatedAt,
		&dto.EmailVerified,
		&dto.SessionVersion,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, localErrors.ErrNotFound
	}

	user := dto.ToModel()

	return &user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	q := `
SELECT id, role_id, username, email, password, created_at, email_verified, session_version 
FROM users 
WHERE email = $1
`

	var dto userDTO

	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&dto.ID,
		&dto.RoleID,
		&dto.Username,
		&dto.Email,
		&dto.Password,
		&dto.CreatedAt,
		&dto.EmailVerified,
		&dto.SessionVersion,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, localErrors.ErrNotFound
	}

	user := dto.ToModel()

	return &user, nil
}

func (r *UserRepo) GetAll(ctx context.Context) ([]model.User, error) {
	q := `select id, role_id, username, email, password, created_at, email_verified, session_version from users`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.User

	for rows.Next() {
		var dto userDTO

		err := rows.Scan(
			&dto.ID,
			&dto.RoleID,
			&dto.Username,
			&dto.Email,
			&dto.Password,
			&dto.CreatedAt,
			&dto.EmailVerified,
			&dto.SessionVersion,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, dto.ToModel())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *UserRepo) Create(ctx context.Context, user model.User) (*model.User, error) {
	q := `INSERT INTO users (role_id, username, email, password, email_verified) 
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, created_at`

	var (
		id        int
		createdAt time.Time
	)

	err := r.db.QueryRowContext(ctx, q,
		user.RoleID,
		user.Username,
		user.Email,
		user.Password,
		user.EmailVerified,
	).Scan(&id, &createdAt)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "users_username_key"):
			return nil, usecase.ErrUserNameNotUnique
		case strings.Contains(err.Error(), "users_email_key"):
			return nil, usecase.ErrUserEmailNotUnique
		default:
			return nil, err
		}
	}

	return &model.User{
		ID:            id,
		RoleID:        user.RoleID,
		Username:      user.Username,
		Email:         user.Email,
		Password:      user.Password,
		CreatedAt:     createdAt,
		EmailVerified: user.EmailVerified,
	}, nil
}

func (r *UserRepo) UpdateEmailVerified(ctx context.Context, email string) error {
	q := `update users set email_verified = true where email = $1 and email_verified = false`

	result, err := r.db.ExecContext(ctx, q, email)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return localErrors.ErrNotFound
	}

	return nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, email, hashedPassword string) error {
	q := `update users set password = $1, session_version = session_version + 1 where email = $2`

	result, err := r.db.ExecContext(ctx, q, hashedPassword, email)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return localErrors.ErrNotFound
	}

	return nil
}

func (r *UserRepo) IncrementSessionVersion(ctx context.Context, email string) error {
	q := `update users set session_version = session_version + 1 where email = $1`

	result, err := r.db.ExecContext(ctx, q, email)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return localErrors.ErrNotFound
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id int) error {
	q := `delete from users where id=$1`

	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return localErrors.ErrNotFound
	}

	return nil
}
