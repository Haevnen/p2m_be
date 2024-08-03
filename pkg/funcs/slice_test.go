package funcs

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestChunk(t *testing.T) {
	t.Run("normal case (integer)", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		chunkSize := 2
		expected := [][]int{{1, 2}, {3, 4}, {5}}
		result := Chunk(input, chunkSize)
		if diff := cmp.Diff(result, expected); diff != "" {
			t.Errorf("expected %v, but got: %v", expected, result)
		}
	})

	t.Run("normal case (string)", func(t *testing.T) {
		input := []string{"A", "B", "C", "D", "E"}
		chunkSize := 2
		expected := [][]string{{"A", "B"}, {"C", "D"}, {"E"}}
		result := Chunk(input, chunkSize)
		if diff := cmp.Diff(result, expected); diff != "" {
			t.Errorf("expected %v, but got: %v", expected, result)
		}
	})

	t.Run("single item", func(t *testing.T) {
		input := []int{1}
		chunkSize := 2
		expected := [][]int{{1}}
		result := Chunk(input, chunkSize)
		if diff := cmp.Diff(result, expected); diff != "" {
			t.Errorf("expected %v, but got: %v", expected, result)
		}
	})

	t.Run("chunk size larger than slice", func(t *testing.T) {
		input := []int{1, 2}
		chunkSize := 5
		expected := [][]int{{1, 2}}
		result := Chunk(input, chunkSize)
		if diff := cmp.Diff(result, expected); diff != "" {
			t.Errorf("expected %v, but got: %v", expected, result)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		chunkSize := 2
		var expected [][]int
		result := Chunk(input, chunkSize)
		if diff := cmp.Diff(result, expected); diff != "" {
			t.Errorf("expected %v, but got: %v", expected, result)
		}
	})

	t.Run("chunk size one", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		chunkSize := 1
		expected := [][]int{{1}, {2}, {3}, {4}, {5}}
		result := Chunk(input, chunkSize)
		if diff := cmp.Diff(result, expected); diff != "" {
			t.Errorf("expected %v, but got: %v", expected, result)
		}
	})
}

func TestContains(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		sl := []string{}
		if Contains(sl, "test") {
			t.Errorf("expected %v, but got: %v", false, true)
		}
	})

	t.Run("contains", func(t *testing.T) {
		sl := []string{"test", "hello", "world"}
		if !Contains(sl, "test") {
			t.Errorf("expected %v, but got: %v", true, false)
		}
	})

	t.Run("not contains", func(t *testing.T) {
		sl := []string{"hello", "world"}
		if Contains(sl, "test") {
			t.Errorf("expected %v, but got: %v", false, true)
		}
	})

	t.Run("contains in int slice", func(t *testing.T) {
		sl := []int{1, 2, 3}
		if !Contains(sl, 2) {
			t.Errorf("expected %v, but got: %v", true, false)
		}
	})

	t.Run("not contains in int slice", func(t *testing.T) {
		sl := []int{1, 2, 3}
		if Contains(sl, 5) {
			t.Errorf("expected %v, but got: %v", false, true)
		}
	})
}

func TestContainErrors(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	err3 := errors.New("error 3")
	errUnknown := errors.New("unknown error")

	t.Run("error present in slice", func(t *testing.T) {
		slice := []error{err1, err2}
		err := err1
		want := true
		if ContainErrors(slice, err) != want {
			t.Errorf("expected %v, but got: %v", want, !want)
		}
	})

	t.Run("error not present in slice", func(t *testing.T) {
		slice := []error{err1, err2}
		err := errUnknown
		want := false
		if ContainErrors(slice, err) != want {
			t.Errorf("expected %v, but got: %v", want, !want)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		slice := []error{}
		err := err3
		want := false
		if ContainErrors(slice, err) != want {
			t.Errorf("expected %v, but got: %v", want, !want)
		}
	})

	t.Run("nil error", func(t *testing.T) {
		slice := []error{err1, err2}
		err := error(nil)
		want := false
		if ContainErrors(slice, err) != want {
			t.Errorf("expected %v, but got: %v", want, !want)
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		slice := []error(nil)
		err := err1
		want := false
		if ContainErrors(slice, err) != want {
			t.Errorf("expected %v, but got: %v", want, !want)
		}
	})
}

func TestRemoveDuplicateItem(t *testing.T) {
	t.Run("integers without duplicates", func(t *testing.T) {
		got := RemoveDuplicateItem([]int{1, 2, 3, 4, 5})
		want := []int{1, 2, 3, 4, 5}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("integers with duplicates", func(t *testing.T) {
		got := RemoveDuplicateItem([]int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5})
		want := []int{1, 2, 3, 4, 5}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("string without duplicates", func(t *testing.T) {
		got := RemoveDuplicateItem([]string{"one", "two", "three", "four", "five"})
		want := []string{"one", "two", "three", "four", "five"}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("string with duplicates", func(t *testing.T) {
		got := RemoveDuplicateItem([]string{"one", "two", "three", "two", "one", "five"})
		want := []string{"one", "two", "three", "five"}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		got := RemoveDuplicateItem([]int{})
		want := []int{}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})
}

func TestFindDuplicate(t *testing.T) {

	t.Run("no duplicates (int)", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		want := false
		got := FindDuplicate(slice)
		if got != want {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("no duplicates (string)", func(t *testing.T) {
		slice := []string{"A", "B", "C"}
		want := false
		got := FindDuplicate(slice)
		if got != want {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("single duplicate (int)", func(t *testing.T) {
		slice := []int{1, 2, 3, 2, 5}
		want := true
		got := FindDuplicate(slice)
		if got != want {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("single duplicate (string)", func(t *testing.T) {
		slice := []string{"A", "B", "C", "B"}
		want := true
		got := FindDuplicate(slice)
		if got != want {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("multiple duplicates", func(t *testing.T) {
		slice := []int{1, 2, 3, 2, 5, 1}
		want := true
		got := FindDuplicate(slice)
		if got != want {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		slice := []int{}
		want := false
		got := FindDuplicate(slice)
		if got != want {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})
}

func TestFindDuplicateVal(t *testing.T) {
	t.Run("no duplicates (int)", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		var want []int
		got := FindDuplicateVal(slice)
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("no duplicates (string)", func(t *testing.T) {
		slice := []string{"A", "B", "C"}
		var want []string
		got := FindDuplicateVal(slice)
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("single duplicate (int)", func(t *testing.T) {
		slice := []int{1, 2, 3, 2, 5}
		want := []int{2}
		got := FindDuplicateVal(slice)
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("single duplicate (string)", func(t *testing.T) {
		slice := []string{"A", "B", "C", "A"}
		want := []string{"A"}
		got := FindDuplicateVal(slice)
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("multiple duplicate", func(t *testing.T) {
		slice := []int{1, 1, 2, 2, 3, 3}
		want := []int{1, 2, 3}
		got := FindDuplicateVal(slice)
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		slice := []int{}
		var want []int
		got := FindDuplicateVal(slice)
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("empty slice (int)", func(t *testing.T) {
		slice := []int{}
		want := []int{}
		got := Filter(slice, func(i int) bool { return i%2 == 0 })
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("empty slice (string)", func(t *testing.T) {
		slice := []string{}
		want := []string{}
		got := Filter(slice, func(s string) bool { return len(s) > 0 })
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("multiple slice (int)", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6}
		want := []int{2, 4, 6}
		got := Filter(slice, func(i int) bool { return i%2 == 0 })
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("multiple slice (string)", func(t *testing.T) {
		slice := []string{"", "A", "B", "", "C"}
		want := []string{"A", "B", "C"}
		got := Filter(slice, func(s string) bool { return len(s) > 0 })
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("empty slice (int to string)", func(t *testing.T) {
		slice := []int{}
		want := []string{}
		got := Map(slice, func(i int) string { return strconv.Itoa(i) })
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("empty slice (string to int)", func(t *testing.T) {
		slice := []string{}
		want := []int{}
		got := Map(slice, func(s string) int {
			i, _ := strconv.Atoi(s)
			return i
		})
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("multiple slice (int to string)", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6}
		want := []string{"1", "2", "3", "4", "5", "6"}
		got := Map(slice, func(i int) string { return strconv.Itoa(i) })
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("multiple slice (string to int)", func(t *testing.T) {
		slice := []string{"1", "2", "3", "4", "5", "6"}
		want := []int{1, 2, 3, 4, 5, 6}
		got := Map(slice, func(s string) int {
			i, _ := strconv.Atoi(s)
			return i
		})
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})
}

func TestFilterMap(t *testing.T) {
	t.Run("empty slice (int to string)", func(t *testing.T) {
		slice := []int{}
		want := []string{}
		got := FilterMap(slice, func(i int) (string, bool) { return strconv.Itoa(i), i%2 == 0 })
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("empty slice (string to int)", func(t *testing.T) {
		slice := []string{}
		want := []int{}
		got := FilterMap(slice, func(s string) (int, bool) {
			i, _ := strconv.Atoi(s)
			return i, i%2 == 0
		})
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("multiple slice (int to string)", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6}
		want := []string{"2", "4", "6"}
		got := FilterMap(slice, func(i int) (string, bool) { return strconv.Itoa(i), i%2 == 0 })
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("multiple slice (string to int)", func(t *testing.T) {
		slice := []string{"1", "2", "3", "4", "5", "6"}
		want := []int{2, 4, 6}
		got := FilterMap(slice, func(s string) (int, bool) {
			i, _ := strconv.Atoi(s)
			return i, i%2 == 0
		})
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})
}

func TestFind(t *testing.T) {
	t.Run("string case", func(t *testing.T) {
		slice := []string{"a", "bb", "ccc", "d", "ee", "fff"}
		want := "bb"
		res := Find(slice, func(s string) bool { return len(s) == 2 })
		if *res != want {
			t.Errorf("expected %v, but got: %v", want, *res)
		}
	})
	t.Run("not found string case", func(t *testing.T) {
		slice := []string{"a", "bb", "ccc", "d", "ee", "fff"}
		res := Find(slice, func(s string) bool { return len(s) == 4 })
		if res != nil {
			t.Errorf("expected nil, but got: %v", res)
		}
	})
	t.Run("int case", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		want := 10
		res := Find(slice, func(i int) bool { return len(strconv.Itoa(i)) == 2 })
		if *res != want {
			t.Errorf("expected %v, but got: %v", want, *res)
		}
	})
	t.Run("struct case", func(t *testing.T) {
		type test struct {
			name string
			age  int
		}
		slice := []*test{{name: "123", age: 123}, {name: "456", age: 456}, {name: "789", age: 789}}
		want := &test{name: "789", age: 789}
		res := Find(slice, func(ts *test) bool { return ts.age > 500 })
		if !reflect.DeepEqual(*res, want) {
			t.Errorf("expected %v, but got: %v", want, *res)
		}
	})
	t.Run("not found struct case", func(t *testing.T) {
		type test struct {
			name string
			age  int
		}
		slice := []*test{{name: "123", age: 123}, {name: "456", age: 456}}
		res := Find(slice, func(ts *test) bool { return ts.age > 500 })
		if res != nil {
			t.Errorf("expected nil, but got: %v", res)
		}
	})
}
