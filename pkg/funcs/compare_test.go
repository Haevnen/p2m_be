package funcs

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/oklog/ulid/v2"
)

func TestMinWithCompare(t *testing.T) {

	t.Run("Compare arrays of strings", func(t *testing.T) {
		slice := []*string{ToPtr("Red"), ToPtr("Orange"), ToPtr("Green")}
		got := MinWithCompare(slice, CompareString)
		want := "Green"
		if *got != want {
			t.Errorf("expected %v, but got: %v", want, *got)
		}
	})

	t.Run("Compare arrays of integers", func(t *testing.T) {
		slice := []*int{ToPtr(5), ToPtr(4), ToPtr(3), ToPtr(2), ToPtr(1)}
		got := MinWithCompare(slice, CompareInt)
		want := 1
		if *got != want {
			t.Errorf("expected %v, but got: %v", want, *got)
		}
	})

	t.Run("Compare arrays of dates", func(t *testing.T) {
		slice := []*string{ToPtr("2022-12-01"), ToPtr("2023-01-01"), ToPtr("2021-11-30")}
		got := MinWithCompare(slice, CompareDate)
		want := "2021-11-30"
		if *got != want {
			t.Errorf("expected %v, but got: %v", want, *got)
		}
	})

	t.Run("Compare arrays of ULIDs", func(t *testing.T) {
		slice := []string{"01FSKYV7SK5F2RWG1WEER2V07D", "01FSKYQJYS6ZB5XFP3N9A7ZWZ8", "01FSKYMM3X4DQV6H410KC60A0M"}
		got := MinWithCompare(slice, CompareUlid)
		want := "01FSKYMM3X4DQV6H410KC60A0M"
		if got != want {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})
}

func TestCompareDate(t *testing.T) {

	targetDate := ToPtr("2023-10-10")
	t.Run("targetDate same dates", func(t *testing.T) {
		result := CompareDate(targetDate, ToPtr("2023-10-10"))
		if result != 0 {
			t.Errorf("expected %v, but got: %v", 0, result)
		}
	})

	t.Run("targetDate after dates", func(t *testing.T) {
		result := CompareDate(targetDate, ToPtr("2023-10-11"))
		if result != -1 {
			t.Errorf("expected %v, but got: %v", -1, result)
		}
	})

	t.Run("targetDate before dates", func(t *testing.T) {
		result := CompareDate(targetDate, ToPtr("2023-10-09"))
		if result != 1 {
			t.Errorf("expected %v, but got: %v", 1, result)
		}
	})

	t.Run("both is nil", func(t *testing.T) {
		result := CompareDate(nil, nil)
		if result != 0 {
			t.Errorf("expected %v, but got: %v", 0, result)
		}
	})

	t.Run("first arg is nil", func(t *testing.T) {
		result := CompareDate(nil, targetDate)
		if result != -1 {
			t.Errorf("expected %v, but got: %v", -1, result)
		}
	})

	t.Run("second arg is nil", func(t *testing.T) {
		result := CompareDate(targetDate, nil)
		if result != 1 {
			t.Errorf("expected %v, but got: %v", 1, result)
		}
	})

	t.Run("equal value", func(t *testing.T) {
		result := CompareDate(targetDate, targetDate)
		if result != 0 {
			t.Errorf("expected %v, but got: %v", 0, result)
		}
	})
}

func TestCompareString(t *testing.T) {
	t.Run("both strings are nil", func(t *testing.T) {
		o := swag.String("")
		target := swag.String("")
		want := 0

		if got := CompareString(o, target); got != want {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("origin string is nil and target string is non-nil", func(t *testing.T) {
		o := swag.String("")
		target := swag.String("test")
		want := 1

		if got := CompareString(o, target); got != want {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("origin string is non-nil and target string is nil", func(t *testing.T) {
		o := swag.String("test")
		target := swag.String("")
		want := -1

		if got := CompareString(o, target); got != want {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("both strings are non-nil and equal", func(t *testing.T) {
		o := swag.String("test")
		target := swag.String("test")
		want := 0

		if got := CompareString(o, target); got != want {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("origin string is 'testA' and target string is 'testB'", func(t *testing.T) {
		o := swag.String("testA")
		target := swag.String("testB")
		want := -1

		if got := CompareString(o, target); got != want {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})

	t.Run("origin string is 'testB' and target string is 'testA'", func(t *testing.T) {
		o := swag.String("testB")
		target := swag.String("testA")
		want := 1

		if got := CompareString(o, target); got != want {
			t.Errorf("expected %v, but got: %v", want, got)
		}
	})
}

func TestCompareInt(t *testing.T) {
	var nilPointer *int

	t.Run("Compare with both nil arguments", func(t *testing.T) {
		if got := CompareInt(nilPointer, nilPointer); got != 0 {
			t.Errorf("expected %v, but got: %v", 0, got)
		}
	})

	t.Run("Compare with first argument nil", func(t *testing.T) {
		arg := swag.Int(5)
		if got := CompareInt(nilPointer, arg); got != -1 {
			t.Errorf("CompareInt() = %v, want -1", got)
		}
	})

	t.Run("Compare with second argument nil", func(t *testing.T) {
		arg := swag.Int(5)
		if got := CompareInt(arg, nilPointer); got != 1 {
			t.Errorf("CompareInt() = %v, want 1", got)
		}
	})

	t.Run("Compare equal non-nil arguments", func(t *testing.T) {
		arg := swag.Int(5)
		if got := CompareInt(arg, arg); got != 0 {
			t.Errorf("expected %v, but got: %v", 0, got)
		}
	})

	t.Run("Compare first argument less than second", func(t *testing.T) {
		arg1 := swag.Int(4)
		arg2 := swag.Int(5)
		if got := CompareInt(arg1, arg2); got != -1 {
			t.Errorf("CompareInt() = %v, want -1", got)
		}
	})

	t.Run("Compare first argument greater than second", func(t *testing.T) {
		arg1 := swag.Int(6)
		arg2 := swag.Int(5)
		if got := CompareInt(arg1, arg2); got != 1 {
			t.Errorf("CompareInt() = %v, want 1", got)
		}
	})
}

func TestCompareULID(t *testing.T) {
	t.Run("both id empty", func(t *testing.T) {
		result := CompareUlid("", "")
		if result != 0 {
			t.Errorf("expected %v, but got: %v", 0, result)
		}
	})

	t.Run("origin empty", func(t *testing.T) {
		targetID := ulid.MustNew(123456789, nil).String()
		result := CompareUlid("", targetID)
		if result != -1 {
			t.Errorf("expected %v, but got: %v", -1, result)
		}
	})

	t.Run("target empty", func(t *testing.T) {
		originID := ulid.MustNew(123456789, nil).String()
		result := CompareUlid(originID, "")
		if result != 1 {
			t.Errorf("expected %v, but got: %v", 1, result)
		}
	})

	t.Run("origin ulid less than target", func(t *testing.T) {
		originID := ulid.MustNew(123456789, nil).String()
		targetID := ulid.MustNew(987654321, nil).String()
		result := CompareUlid(originID, targetID)
		if result != -1 {
			t.Errorf("expected %v, but got: %v", -1, result)
		}
	})

	t.Run("origin ulid greater than target", func(t *testing.T) {
		originID := ulid.MustNew(987654321, nil).String()
		targetID := ulid.MustNew(123456789, nil).String()
		result := CompareUlid(originID, targetID)
		if result != 1 {
			t.Errorf("expected %v, but got: %v", 1, result)
		}
	})

	t.Run("origin ulid equals target", func(t *testing.T) {
		id := ulid.MustNew(123456789, nil).String()
		result := CompareUlid(id, id)
		if result != 0 {
			t.Errorf("expected %v, but got: %v", 0, result)
		}
	})
}

func TestCompareNil(t *testing.T) {
	var argStr = ToPtr("test")
	var argNil *string = nil

	t.Run("both nil and isDesc is false", func(t *testing.T) {
		got, isNil := CompareNil(argNil, argNil, false)
		if !isNil {
			t.Errorf("expected %v, but got: %v", true, isNil)
		}
		if got != 0 {
			t.Errorf("expected %v, but got: %v", 0, got)
		}
	})

	t.Run("both nil and isDesc is true", func(t *testing.T) {
		got, isNil := CompareNil(argNil, argNil, true)
		if !isNil {
			t.Errorf("expected %v, but got: %v", true, isNil)
		}
		if got != 0 {
			t.Errorf("expected %v, but got: %v", 0, got)
		}
	})

	t.Run("both not nil and isDesc is false", func(t *testing.T) {
		got, isNil := CompareNil(argStr, argStr, false)
		if isNil {
			t.Errorf("expected %v, but got: %v", false, isNil)
		}
		if got != 0 {
			t.Errorf("expected %v, but got: %v", 0, got)
		}
	})

	t.Run("both not nil and isDesc is true", func(t *testing.T) {
		got, isNil := CompareNil(argStr, argStr, true)
		if isNil {
			t.Errorf("expected %v, but got: %v", false, isNil)
		}
		if got != 0 {
			t.Errorf("expected %v, but got: %v", 0, got)
		}
	})

	t.Run("origin nil and isDesc is false", func(t *testing.T) {
		got, isNil := CompareNil(argNil, argStr, false)
		if !isNil {
			t.Errorf("expected %v, but got: %v", true, isNil)
		}
		if got != 1 {
			t.Errorf("expected %v, but got: %v", 1, got)
		}
	})

	t.Run("origin nil and isDesc is true", func(t *testing.T) {
		got, isNil := CompareNil(argNil, argStr, true)
		if !isNil {
			t.Errorf("expected %v, but got: %v", true, isNil)
		}
		if got != -1 {
			t.Errorf("expected %v, but got: %v", -1, got)
		}
	})

	t.Run("target nil and isDesc is false", func(t *testing.T) {
		got, isNil := CompareNil(argStr, argNil, false)
		if !isNil {
			t.Errorf("expected %v, but got: %v", true, isNil)
		}
		if got != -1 {
			t.Errorf("expected %v, but got: %v", -1, got)
		}
	})

	t.Run("target nil and isDesc is true", func(t *testing.T) {
		got, isNil := CompareNil(argStr, argNil, true)
		if !isNil {
			t.Errorf("expected %v, but got: %v", true, isNil)
		}
		if got != 1 {
			t.Errorf("expected %v, but got: %v", 1, got)
		}
	})
}
