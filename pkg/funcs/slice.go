package funcs

import (
	"errors"
	"slices"
)

// Chunk chunk slice
func Chunk[T any](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	for {
		if len(slice) == 0 {
			break
		}
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}
		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}
	return chunks
}

// Contains is contain in slice
func Contains[T comparable](slice []T, v T) bool {
	for _, s := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// ContainErrors is contain error in slice
func ContainErrors(slice []error, err error) bool {
	for _, e := range slice {
		if errors.Is(err, e) {
			return true
		}
	}
	return false
}

// RemoveDuplicateItem is remove duplicate item
func RemoveDuplicateItem[T comparable](slice []T) []T {
	results := make([]T, 0, len(slice))
	encountered := map[T]bool{}
	for _, s := range slice {
		if !encountered[s] {
			encountered[s] = true
			results = append(results, s)
		}
	}
	return slices.Clip(results)
}

// FindDuplicate return true when slice has duplicate value
func FindDuplicate[T comparable](slice []T) bool {
	uniqMap := make(map[T]bool)
	for _, v := range slice {
		if _, ok := uniqMap[v]; !ok {
			uniqMap[v] = true
		} else {
			return true
		}
	}
	return false
}

// FindDuplicateVal returns slice of duplicate values when slice has duplicate value
func FindDuplicateVal[T comparable](slice []T) []T {
	uniqMap := make(map[T]bool)
	var duplicates []T
	for _, v := range slice {
		if _, ok := uniqMap[v]; !ok {
			uniqMap[v] = true
		} else {
			duplicates = append(duplicates, v)
		}
	}
	return slices.Clip(duplicates)
}

// Filter filter slice
func Filter[T any](slice []T, fn func(T) bool) []T {
	outputs := make([]T, 0)
	for _, s := range slice {
		if fn(s) {
			outputs = append(outputs, s)
		}
	}
	return slices.Clip(outputs)
}

// Map map slice
func Map[T, V any](slice []T, fn func(T) V) []V {
	outputs := make([]V, len(slice))
	for i, s := range slice {
		outputs[i] = fn(s)
	}
	return outputs
}

// FilterMap filter and map slice
func FilterMap[T, V any](slice []T, fn func(T) (V, bool)) []V {
	outputs := make([]V, 0, len(slice))
	for _, s := range slice {
		v, ok := fn(s)
		if ok {
			outputs = append(outputs, v)
		}
	}
	return slices.Clip(outputs)
}

// Find find elements in slice
func Find[T any](slice []T, fn func(T) bool) *T {
	for _, s := range slice {
		if fn(s) {
			return &s
		}
	}
	return nil
}
