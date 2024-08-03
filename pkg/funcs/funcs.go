// Package funcs define general functions
package funcs

// ToPtr convert pointer
func ToPtr[T any](t T) *T {
	return &t
}

// ToPtrValue convert value in pointer
func ToPtrValue[T any](t *T) T {
	if t == nil {
		out := (*T)(nil)
		return *out
	}
	return *t
}

// ToPtrValueSafe safe convert value in pointer
func ToPtrValueSafe[T any](t *T) T {
	if t == nil {
		return *new(T) // nolint:gocritic
	}
	return *t
}

// IsNil check is nil
func IsNil[T any](t *T) bool {
	return t == nil
}

// IsNilOrZero check is nil or init value
func IsNilOrZero[T comparable](t *T) bool {
	if t == nil {
		return true
	}
	return *t == *new(T) // nolint:gocritic
}
