package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
	"palback/internal/pkg/helpers"
	tokens "palback/internal/pkg/token"
	ucModel "palback/internal/usecase/model"
	"palback/internal/usecase/port"
)

type UserUseCase struct {
	roleService RoleService
	mailer      port.EmailSender
	kvStorage   port.KeyValueStorage
	repo        port.UserRepo
}

func NewUserUseCase(
	roleService RoleService,
	mailer port.EmailSender,
	kvStorage port.KeyValueStorage,
	repo port.UserRepo,
) *UserUseCase {
	return &UserUseCase{
		roleService: roleService,
		mailer:      mailer,
		kvStorage:   kvStorage,
		repo:        repo,
	}
}

func (s *UserUseCase) Get(ctx context.Context, id int) (*ucModel.UserDetail, error) {
	user, err := s.repo.Get(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя по id: %w", err)
	}

	role, err := s.roleService.Get(ctx, user.RoleID)

	if err != nil {
		return nil, fmt.Errorf("ошибка получения роли по id: %w", err)
	}

	userDetail := ucModel.CreateUserDetail(helpers.FromPtr(user), helpers.FromPtr(role))

	return &userDetail, nil
}

func (s *UserUseCase) GetAll(ctx context.Context) (result ucModel.UserList, err error) {
	users, err := s.repo.GetAll(ctx)

	if err != nil {
		return result, fmt.Errorf("ошибка при получении списка регионов: %w", err)
	}

	rolesMap, err := s.roleService.GetAllMap(ctx)
	if err != nil {
		return result, fmt.Errorf("ошибка при получении информации о ролях: %w", err)
	}

	result = ucModel.CreateUserList(users, rolesMap)

	return result, nil
}

func (s *UserUseCase) Register(
	ctx context.Context,
	userName, email, password string,
) (_ *ucModel.UserDetail, err error) {
	var (
		user  *model.User
		token string
	)

	defer func() {
		// Если возникли проблемы с регистрацией, то созданная запись удаляется
		if err != nil && user != nil {
			err = s.repo.Delete(ctx, user.ID)
			err = s.kvStorage.Del(ctx, "verify_email:"+token)
		}
	}()

	// Сгенерировать хэш пароля
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации пароля: %w", err)
	}

	// Добавить пользователя
	user, err = s.repo.Create(ctx, model.User{
		RoleID:        model.RoleUser,
		Username:      userName,
		Email:         email,
		Password:      string(hashed),
		EmailVerified: false,
	})

	if err != nil {
		switch {
		default:
			return nil, fmt.Errorf("ошибка добавления пользователя: %w", err)
		}
	}

	// Создать токен для проверки и поместить его в key-value хранилище
	token, err = tokens.GenerateVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации токена: %w", err)
	}

	err = s.kvStorage.Set(ctx, "verify_email:"+token, email, 3600)
	if err != nil {
		return nil, fmt.Errorf("ошибка записи токена в хранилище: %w", err)
	}

	// Создать и отправить проверочное письмо
	err = s.mailer.SendVerificationEmail(email, token)
	if err != nil {
		log.Printf("Ошибка отправки проверочного письма на %s: %v", email, err)
		return nil, fmt.Errorf("ошибка отправки проверочного письма: %w", err)
	}

	// Собрать детальную информацию о пользователе
	role, err := s.roleService.Get(ctx, user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения роли по id: %w", err)
	}

	userDetail := ucModel.CreateUserDetail(helpers.FromPtr(user), helpers.FromPtr(role))

	return &userDetail, nil
}

func (s *UserUseCase) VerifyEmail(ctx context.Context, token string) error {
	// Получить значение токена
	email, err := s.kvStorage.Get(ctx, "verify_email:"+token)
	if err != nil {
		switch {
		case errors.Is(err, ErrNoReplyFromKeyValueStorage):
			return fmt.Errorf("неверный или устаревший токен: %w", err)
		default:
			return fmt.Errorf("внутренняя ошибка получения токена: %w", err)
		}
	}

	// Обновить информацию о пользователе
	err = s.repo.UpdateEmailVerified(ctx, email)
	if err != nil {
		return err
	}

	// Удалить токен
	_ = s.kvStorage.Del(ctx, "verify_email:"+token)

	return nil
}

func (s *UserUseCase) ResendVerificationEmail(ctx context.Context, email string) error {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil
	}

	if user.EmailVerified {
		return nil
	}

	tokenStr, err := tokens.GenerateVerificationToken()
	if err != nil {
		return fmt.Errorf("ошибка генерации токена: %w", err)
	}

	err = s.kvStorage.Set(ctx, "verify_email:"+tokenStr, email, 3600)
	if err != nil {
		return fmt.Errorf("ошибка сохранения токена: %w", err)
	}

	err = s.mailer.SendVerificationEmail(email, tokenStr)
	if err != nil {
		log.Println("ошибка отправки письма пользователю", err)
	}

	return nil

}

func (s *UserUseCase) Login(ctx context.Context, identifier, password string) (*ucModel.UserDetail, error) {
	user, err := s.repo.GetByIdentifier(ctx, identifier)
	if err != nil || user == nil {
		switch {
		case errors.Is(err, localErrors.ErrNotFound):
			return nil, ErrUserInvalidCredentials
		default:
			return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrUserInvalidCredentials
	}

	if !user.EmailVerified {
		return nil, ErrUncheckedEmail
	}

	// Собрать детальную информацию о пользователе
	role, err := s.roleService.Get(ctx, user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения роли по id: %w", err)
	}

	userDetail := ucModel.CreateUserDetail(helpers.FromPtr(user), helpers.FromPtr(role))

	return &userDetail, nil
}

func (s *UserUseCase) RequestPasswordReset(ctx context.Context, email string) error {
	// Проверяем, существует ли пользователь
	_, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		// Не раскрываем, существует ли email — для безопасности
		return nil
	}

	// Генерируем токен
	tokenStr, err := tokens.GenerateVerificationToken()
	if err != nil {
		return fmt.Errorf("ошибка генерации токена: %w", err)
	}

	// Сохраняем токен в хранилище
	err = s.kvStorage.Set(ctx, "verify_email:"+tokenStr, email, 3600)
	if err != nil {
		return fmt.Errorf("ошибка сохранения токена: %w", err)
	}

	err = s.mailer.SendPasswordResetEmail(email, tokenStr)
	if err != nil {
		log.Println("ошибка отправки письма пользователю", err)
	}

	return nil
}

func (s *UserUseCase) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	if token == "" || newPassword == "" {
		return ErrInvalidToken
	}

	// получаем email по токену
	email, err := s.kvStorage.Get(ctx, "reset_token:"+token)
	if err != nil {
		switch {
		case errors.Is(err, ErrNoReplyFromKeyValueStorage):
			return ErrInvalidToken
		default:
			return fmt.Errorf("ошибка получения токена: %w", err)
		}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.repo.UpdatePassword(ctx, email, string(hashed))
	if err != nil {
		return fmt.Errorf("ошибка обновления пароля: %w", err)
	}

	_ = s.kvStorage.Del(ctx, "reset_token:"+token)

	return nil
}
