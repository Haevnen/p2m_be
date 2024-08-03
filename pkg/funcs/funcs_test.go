package funcs

import "testing"

func TestToPtr(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		input := "test"
		got := ToPtr(input)
		if got == nil {
			t.Errorf("expected not nil, but got: %v", got)
		}
		if *got != input {
			t.Errorf("expected %v, but got: %v", input, got)
		}
	})

	t.Run("int", func(t *testing.T) {
		input := 1
		got := ToPtr(input)
		if got == nil {
			t.Errorf("expected not nil, but got: %v", got)
		}
		if *got != input {
			t.Errorf("expected %v, but got: %v", input, got)
		}
	})

	t.Run("bool", func(t *testing.T) {
		input := true
		got := ToPtr(input)
		if got == nil {
			t.Errorf("expected not nil, but got: %v", got)
		}
		if *got != input {
			t.Errorf("expected %v, but got: %v", input, got)
		}
	})

	t.Run("empty string", func(t *testing.T) {
		input := ""
		got := ToPtr(input)
		if got == nil {
			t.Errorf("expected not nil, but got: %v", got)
		}
		if *got != input {
			t.Errorf("expected %v, but got: %v", input, got)
		}
	})

	t.Run("zero value int", func(t *testing.T) {
		input := 0
		got := ToPtr(input)
		if got == nil {
			t.Errorf("expected not nil, but got: %v", got)
		}
		if *got != input {
			t.Errorf("expected %v, but got: %v", input, got)
		}
	})

	t.Run("false boolean", func(t *testing.T) {
		input := false
		got := ToPtr(input)
		if got == nil {
			t.Errorf("expected not nil, but got: %v", got)
		}
		if *got != input {
			t.Errorf("expected %v, but got: %v", input, got)
		}
	})
}

func TestToPtrValue(t *testing.T) {
	t.Run("string nil value", func(t *testing.T) {
		var input *string = nil
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected nil pointer dereference, but no error occurred")
			} else {
				if r.(error).Error() != "runtime error: invalid memory address or nil pointer dereference" {
					t.Errorf("expected nil pointer dereference, but %v", r)
				}
			}
		}()
		ToPtrValue(input)
	})

	t.Run("string non nil value", func(t *testing.T) {
		val := "ABC"
		input := &val
		want := "ABC"
		if got := ToPtrValue(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})

	t.Run("int nil value", func(t *testing.T) {
		var input *int = nil
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected nil pointer dereference, but no error occurred")
			} else {
				if r.(error).Error() != "runtime error: invalid memory address or nil pointer dereference" {
					t.Errorf("expected nil pointer dereference, but %v", r)
				}
			}
		}()
		ToPtrValue(input)
	})

	t.Run("int non nil value", func(t *testing.T) {
		val := 1
		input := &val
		want := 1
		if got := ToPtrValue(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})
}

func TestToPtrValueSafe(t *testing.T) {
	t.Run("string nil value", func(t *testing.T) {
		var input *string = nil
		want := ""
		if got := ToPtrValueSafe(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})

	t.Run("string non nil value", func(t *testing.T) {
		val := "ABC"
		input := &val
		want := "ABC"
		if got := ToPtrValueSafe(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})

	t.Run("int nil value", func(t *testing.T) {
		var input *int = nil
		want := 0
		if got := ToPtrValueSafe(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})

	t.Run("int non nil value", func(t *testing.T) {
		val := 1
		input := &val
		want := 1
		if got := ToPtrValueSafe(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})
}

func TestIsNil(t *testing.T) {
	t.Run("string nil", func(t *testing.T) {
		var input *string = nil
		want := true
		if got := IsNil(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})

	t.Run("string zero", func(t *testing.T) {
		input := ToPtr("")
		want := false
		if got := IsNil(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})

	t.Run("string non zero", func(t *testing.T) {
		input := ToPtr("ABC")
		want := false
		if got := IsNil(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})

	t.Run("int nil", func(t *testing.T) {
		var input *int = nil
		want := true
		if got := IsNil(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})

	t.Run("int zero", func(t *testing.T) {
		input := ToPtr(0)
		want := false
		if got := IsNil(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})

	t.Run("int non zero", func(t *testing.T) {
		input := ToPtr("123")
		want := false
		if got := IsNil(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})
}

func TestIsNilOrZero(t *testing.T) {
	t.Run("string nil", func(t *testing.T) {
		var input *string = nil
		want := true
		if got := IsNilOrZero(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})

	t.Run("string zero", func(t *testing.T) {
		input := ToPtr("")
		want := true
		if got := IsNilOrZero(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})

	t.Run("string non zero", func(t *testing.T) {
		input := ToPtr("ABC")
		want := false
		if got := IsNilOrZero(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})

	t.Run("int nil", func(t *testing.T) {
		var input *int = nil
		want := true
		if got := IsNilOrZero(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})

	t.Run("int zero", func(t *testing.T) {
		input := ToPtr(0)
		want := true
		if got := IsNilOrZero(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})

	t.Run("int non zero", func(t *testing.T) {
		input := ToPtr("123")
		want := false
		if got := IsNilOrZero(input); got != want {
			t.Errorf("expected not nil, but got: %v", got)
		}
	})
}
