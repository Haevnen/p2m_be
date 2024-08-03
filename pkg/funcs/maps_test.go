package funcs

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMapValues(t *testing.T) {

	t.Run("with not empty map", func(t *testing.T) {
		input := map[string]int{
			"first":  1,
			"second": 2,
			"third":  3,
		}
		expected := []int{1, 2, 3}

		result := MapValues(input)
		sort.Ints(result)

		if diff := cmp.Diff(result, expected); diff != "" {
			t.Errorf("expected %v, but got: %v", expected, result)
		}
	})

	t.Run("with empty map", func(t *testing.T) {
		input := map[string]int{}
		expected := []int{}

		result := MapValues(input)

		if diff := cmp.Diff(result, expected); diff != "" {
			t.Errorf("expected %v, but got: %v", expected, result)
		}
	})

	t.Run("with nil map", func(t *testing.T) {
		var input map[string]int
		expected := []int{}

		result := MapValues(input)

		if diff := cmp.Diff(result, expected); diff != "" {
			t.Errorf("expected %v, but got: %v", expected, result)
		}
	})

	t.Run("with custom type", func(t *testing.T) {
		type CustomMapType map[string]int

		input := CustomMapType{
			"first":  1,
			"second": 2,
			"third":  3,
		}

		expected := []int{1, 2, 3}

		result := MapValues(input)
		sort.Ints(result)

		if diff := cmp.Diff(result, expected); diff != "" {
			t.Errorf("expected %v, but got: %v", expected, result)
		}
	})

}
