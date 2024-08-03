package funcs

import (
	"reflect"
	"strings"
	"time"

	"github.com/go-openapi/swag"
	"github.com/oklog/ulid/v2"

	"github.com/Haevnen/p2m_be/pkg/dateutil"
)

// MinWithCompare extract max value from slice
func MinWithCompare[T any](slice []T, f func(T, T) int) T {
	var minValue T
	for i, v := range slice {
		if i == 0 {
			minValue = v
			continue
		}
		compareResult := f(minValue, v)
		if compareResult > 0 {
			minValue = v
		}
	}
	return minValue
}

// CompareDate compareDate
func CompareDate(origin *string, target *string) int {
	originTime, _ := time.Parse(dateutil.LayoutDate, swag.StringValue(origin)) //nolint:errcheck
	targetTime, _ := time.Parse(dateutil.LayoutDate, swag.StringValue(target)) //nolint:errcheck
	if originTime.After(targetTime) {
		return 1
	}
	if originTime.Before(targetTime) {
		return -1
	}
	return 0
}

// CompareString compareString
func CompareString(originString *string, targetString *string) int {
	if swag.StringValue(originString) == "" && swag.StringValue(targetString) == "" {
		return 0
	}
	if swag.StringValue(originString) == "" {
		return 1
	}
	if swag.StringValue(targetString) == "" {
		return -1
	}
	return strings.Compare(swag.StringValue(originString), swag.StringValue(targetString))
}

// CompareInt compareInt
func CompareInt(originInt *int, targetInt *int) int {
	if swag.IntValue(originInt) == swag.IntValue(targetInt) {
		return 0
	}
	if swag.IntValue(originInt) < swag.IntValue(targetInt) {
		return -1
	}
	return 1
}

// CompareUlid compare ulid
func CompareUlid(originID, targetID string) int {
	if originID == "" && targetID == "" {
		return 0
	}
	if originID == "" {
		return -1
	}
	if targetID == "" {
		return 1
	}
	originUlid := ulid.MustParseStrict(originID)
	targetUlid := ulid.MustParseStrict(targetID)
	return originUlid.Compare(targetUlid)
}

// CompareNil call when either of the values to compare is nil
func CompareNil(origin, target interface{}, isDesc bool) (int, bool) {

	// nil is always at the bottom
	multiply := 1
	if isDesc {
		multiply = -1
	}

	if reflect.ValueOf(origin).IsNil() && reflect.ValueOf(target).IsNil() {
		return 0, true
	}
	if reflect.ValueOf(origin).IsNil() {
		return 1 * multiply, true
	}
	if reflect.ValueOf(target).IsNil() {
		return -1 * multiply, true
	}
	return 0, false
}
