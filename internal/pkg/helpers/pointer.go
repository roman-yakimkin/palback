package helpers

func FromPtr[T any](src *T) (result T) {
	if src != nil {
		result = *src
	}

	return result
}
