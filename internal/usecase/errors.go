package usecase

import "errors"

var (
	ErrCountryNotFound      = errors.New("страна не найдена")
	ErrCountryAlreadyAdded  = errors.New("страна с таким id уже добавлена")
	ErrCountryNameNotUnique = errors.New("страна должна иметь уникальное название")

	ErrCountryHasNotRegions = errors.New("страна не имеет регионов")
	ErrRegionNotFound       = errors.New("регион не найден")
	ErrRegionNotUnique      = errors.New("в пределах одной страны регион должен иметь уникальное название")

	ErrCityTypeNotFound = errors.New("тип населенного пункта не найден")

	ErrPlaceTypeNotFound = errors.New("тип святого места не найден")

	ErrUserNameNotUnique           = errors.New("имя пользователя должно быть уникальным")
	ErrUserEmailNotUnique          = errors.New("e-mail пользователя должен быть уникальным")
	ErrVerificationEmailSendFailed = errors.New("пользователь создан, но проверочное письмо отправить не удалось")
	ErrUserInvalidCredentials      = errors.New("неверные логин и пароль")
	ErrUnauthenticated             = errors.New("не аутентифицировано")
	ErrUncheckedEmail              = errors.New("ваш e-mail должен быть подтвержден, проверьте почту и подвердите e-mail")
	ErrInvalidToken                = errors.New("неверный или устаревший токен")
	ErrSessionExpired              = errors.New("сессия устарела, требуется повторный вход на сайт")

	ErrNoReplyFromKeyValueStorage = errors.New("нет ответа от key-value хранилища")
	ErrKeyNotFound                = errors.New("ключ не найден")
)
