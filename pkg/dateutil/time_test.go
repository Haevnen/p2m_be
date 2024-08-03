package dateutil

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/google/go-cmp/cmp"
	"github.com/oapi-codegen/runtime/types"
	"google.golang.org/genproto/googleapis/type/date"
)

func TestNow(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	TimeNowFunc = func() time.Time { return testTime }

	res := Now()
	if res != testTime {
		t.Errorf("expected %v, but got: %v", testTime, res)
	}
}

func TestJST(t *testing.T) {
	res := JST()
	expect := "Asia/Tokyo"
	if res.String() != expect {
		t.Errorf("expected %v, but got: %v", expect, res)
	}
}

func TestTransactionDefaultDate(t *testing.T) {
	res := TransactionDefaultDate()
	expect := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	if res != expect {
		t.Errorf("expected %v, but got: %v", expect, res)
	}
}

func TestTransactionDefaultDatePtr(t *testing.T) {
	res := TransactionDefaultDatePtr()
	expect := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	if *res != expect {
		t.Errorf("expected %v, but got: %v", expect, *res)
	}
}

func TestValidDefaultDate(t *testing.T) {
	res := ValidDefaultDate()
	expect := time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
	if res != expect {
		t.Errorf("expected %v, but got: %v", expect, res)
	}
}

func TestValidFromDefaultDate(t *testing.T) {
	res := ValidFromDefaultDate()
	expect := time.Date(1000, 01, 01, 0, 0, 0, 0, time.UTC)
	if res != expect {
		t.Errorf("expected %v, but got: %v", expect, res)
	}
}

func TestValidTargetToday(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)
	TimeNowFunc = func() time.Time { return testTime }

	res := ValidTargetToday()
	expect := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if res != expect {
		t.Errorf("expected %v, but got: %v", expect, res)
	}
}

func TestValidDefaultDatePtr(t *testing.T) {
	res := ValidDefaultDatePtr()
	expect := time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
	if *res != expect {
		t.Errorf("expected %v, but got: %v", expect, *res)
	}
}

func TestDateToString(t *testing.T) {
	t.Run("arg is nil", func(t *testing.T) {
		res := DateToString(nil)
		if res != "" {
			t.Errorf("expected %v, but got: %v", "", res)
		}
	})
	t.Run("arg is not nil", func(t *testing.T) {
		arg := types.Date{Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)}
		res := DateToString(&arg)
		expect := "2023-01-01"
		if res != expect {
			t.Errorf("expected %v, but got: %v", expect, res)
		}
	})
}

func TestDateToTime(t *testing.T) {
	t.Run("arg is nil", func(t *testing.T) {
		res := DateToTime(nil)
		if res != nil {
			t.Errorf("expected %v, but got: %v", nil, res)
		}
	})
	t.Run("arg is not nil", func(t *testing.T) {
		arg := types.Date{Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)}
		res := DateToTime(&arg)
		expect := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		if *res != expect {
			t.Errorf("expected %v, but got: %v", expect, *res)
		}
	})
}

func TestTimeToString(t *testing.T) {
	arg := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	res := TimeToString(arg)
	expect := "2023-01-01 00:00:00"
	if res != expect {
		t.Errorf("expected %v, but got: %v", expect, res)
	}
}

func TestTimeToDateString(t *testing.T) {
	arg := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	res := TimeToDateString(arg)
	expect := "2023-01-01"
	if res != expect {
		t.Errorf("expected %v, but got: %v", expect, res)
	}
}

func TestTimeValueDefault(t *testing.T) {
	t.Run("arg is nil", func(t *testing.T) {
		testTime := time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)
		TimeNowFunc = func() time.Time { return testTime }
		res := TimeValueDefault(nil)
		if res != testTime {
			t.Errorf("expected %v, but got: %v", testTime, res)
		}
	})
	t.Run("arg is not nil", func(t *testing.T) {
		testTime := time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)
		TimeNowFunc = func() time.Time { return testTime }

		arg := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		res := TimeValueDefault(&arg)
		if res != arg {
			t.Errorf("expected %v, but got: %v", arg, res)
		}
	})
}

func TestValidFromDefaultTypesDate(t *testing.T) {
	res := ValidFromDefaultTypesDate()
	expect := types.Date{Time: time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC)}
	if res != expect {
		t.Errorf("expected %v, but got: %v", expect, res)
	}
}

func TestIsBeforeOfficialMinDateTime(t *testing.T) {
	t.Run("min date", func(t *testing.T) {
		arg := time.Date(0001, 1, 1, 0, 0, 0, 0, time.UTC)
		res := IsBeforeOfficialMinDateTime(arg)
		if !res {
			t.Errorf("expected %v, but got: %v", true, res)
		}
	})
	t.Run("not min date", func(t *testing.T) {
		arg := time.Date(1001, 1, 1, 0, 0, 0, 0, time.UTC)
		res := IsBeforeOfficialMinDateTime(arg)
		if res {
			t.Errorf("expected %v, but got: %v", false, res)
		}
	})
}

func TestPbDateToDate(t *testing.T) {
	res := PbDateToDate(&date.Date{Year: 2023, Month: 1, Day: 1})
	expect := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if res != expect {
		t.Errorf("expected %v, but got: %v", expect, res)
	}
}

func TestTimeToPbDate(t *testing.T) {
	arg := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	res := TimeToPbDate(arg)
	expect := &date.Date{Year: 2023, Month: 1, Day: 1}
	// this structure cannot use cmp.Diff
	if !reflect.DeepEqual(res, expect) {
		t.Errorf("expected %v, but got: %v", expect, res)
	}
}

func TestBetween(t *testing.T) {
	baseAt := time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC)
	cases := map[string]struct {
		base   time.Time
		from   time.Time
		to     time.Time
		expect bool
	}{
		"from > base < to": {expect: false, base: baseAt, from: baseAt.Add(time.Microsecond), to: baseAt.Add(time.Microsecond)},
		"from = base < to": {expect: false, base: baseAt, from: baseAt, to: baseAt.Add(time.Microsecond)},
		"from < base > to": {expect: false, base: baseAt, from: baseAt.Add(-1 * time.Microsecond), to: baseAt.Add(-1 * time.Microsecond)},
		"from < base = to": {expect: false, base: baseAt, from: baseAt.Add(-1 * time.Microsecond), to: baseAt},
		"from < base < to": {expect: true, base: baseAt, from: baseAt.Add(-1 * time.Microsecond), to: baseAt.Add(time.Microsecond)},
	}
	for caseName, cc := range cases {
		t.Run(caseName, func(t *testing.T) {
			res := Between(cc.base, cc.from, cc.to)
			if res != cc.expect {
				t.Errorf("expected %v, but got: %v", cc.expect, res)
			}
		})
	}
}

func TestBetweenEqual(t *testing.T) {
	baseAt := time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC)
	cases := map[string]struct {
		base   time.Time
		from   time.Time
		to     time.Time
		expect bool
	}{
		"from > base < to": {expect: false, base: baseAt, from: baseAt.Add(time.Microsecond), to: baseAt.Add(time.Microsecond)},
		"from = base < to": {expect: true, base: baseAt, from: baseAt, to: baseAt.Add(time.Microsecond)},
		"from < base > to": {expect: false, base: baseAt, from: baseAt.Add(-1 * time.Microsecond), to: baseAt.Add(-1 * time.Microsecond)},
		"from < base = to": {expect: true, base: baseAt, from: baseAt.Add(-1 * time.Microsecond), to: baseAt},
		"from < base < to": {expect: true, base: baseAt, from: baseAt.Add(-1 * time.Microsecond), to: baseAt.Add(time.Microsecond)},
	}
	for caseName, cc := range cases {
		t.Run(caseName, func(t *testing.T) {
			res := BetweenEqual(cc.base, cc.from, cc.to)
			if res != cc.expect {
				t.Errorf("expected %v, but got: %v", cc.expect, res)
			}
		})
	}
}

func TestDateUnmarshalJSON(t *testing.T) {

	t.Run("json syntax error", func(t *testing.T) {
		n := &Date{}
		err := n.UnmarshalJSON([]byte("0000-00-00"))
		var syntaxError *json.SyntaxError
		if !errors.As(err, &syntaxError) {
			t.Errorf("expected err %v, but got: %v", syntaxError, err)
		}
	})

	t.Run("time parse error", func(t *testing.T) {
		d := &Date{}
		err := d.UnmarshalJSON([]byte("\"2023:01:01\""))
		var parseError *time.ParseError
		if !errors.As(err, &parseError) {
			t.Errorf("expected err %v, but got: %v", parseError, err)
		}
	})

	t.Run("valid time format", func(t *testing.T) {
		d := &Date{}
		err := d.UnmarshalJSON([]byte("\"2023-01-01\""))
		if err != nil {
			t.Errorf("expected err nil, but got: %v", err)
		}

		expect := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		if d.Time != expect {
			t.Errorf("expected %v, but got: %v", expect, d.Time)
		}
	})
}

func TestDateMarshalJSON(t *testing.T) {
	d := &Date{Format: LayoutDate, Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)}
	res, err := d.MarshalJSON()
	if err != nil {
		t.Errorf("expected err nil, but got: %v", err)
	}

	expect := []byte("\"2023-01-01\"")
	if diff := cmp.Diff(res, expect); diff != "" {
		t.Errorf("expected %v, but got: %v", expect, res)
	}
}

func TestStringDateToTime(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "success: 1",
			date:    "2023/01/01",
			want:    time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "success: 2",
			date:    "2023/1/1",
			want:    time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "success: 3",
			date:    "2023/1/01",
			want:    time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "success: 4",
			date:    "2023/01/1",
			want:    time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "failure: 1. invalid format",
			date:    "2023-01-01",
			want:    time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC),
			wantErr: true,
		},
		{
			name:    "failure: 2. The date that do not exist (2023 is not leap year)",
			date:    "2023/02/29",
			want:    time.Date(2023, 02, 29, 0, 0, 0, 0, time.UTC),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringDateToTime(tt.date)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ChooseFormatDate() got = %v, want %v", got, tt.want)
				}
				return
			} else {
				if err != nil || (swag.TimeValue(got) != tt.want) {
					t.Errorf("ChooseFormatDate() got = %v, want %v", got, tt.want)
				}
				return
			}
		})
	}
}
