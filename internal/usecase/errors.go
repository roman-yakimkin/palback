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
)
