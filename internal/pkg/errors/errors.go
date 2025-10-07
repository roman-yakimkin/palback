package errors

import "errors"

var ErrNotFound = errors.New("запись не найдена")

func IsOneOf(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
