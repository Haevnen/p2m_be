// Package dateutil define time functions
package dateutil

import (
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"github.com/go-openapi/swag"
	"github.com/oapi-codegen/runtime/types"
	"google.golang.org/genproto/googleapis/type/date"
)

const (
	// LayoutDate valid layout for date
	LayoutDate = "2006-01-02"

	// LayoutDateTime valid layout for date time
	LayoutDateTime = "2006-01-02 15:04:05"
)

// TimeNowFunc Function to return the current time
var TimeNowFunc = time.Now

// Now Function to make TimeNowFunc() testable
func Now() time.Time {
	return TimeNowFunc()
}

// JST timezone for jst
func JST() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}

// TransactionDefaultDate is default date for BiTemporal
func TransactionDefaultDate() time.Time {
	return time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
}

// TransactionDefaultDatePtr is default date for BiTemporal
func TransactionDefaultDatePtr() *time.Time {
	return swag.Time(TransactionDefaultDate())
}

// ValidDefaultDate is default date for BiTemporal
func ValidDefaultDate() time.Time {
	return time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
}

// ValidFromDefaultDate is default valid_from for transaction model
func ValidFromDefaultDate() time.Time {
	return time.Date(1000, 01, 01, 0, 0, 0, 0, time.UTC)
}

// ValidTargetToday is today for compare valid date
func ValidTargetToday() time.Time {
	now := Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

// ValidDefaultDatePtr is default date for BiTemporal
func ValidDefaultDatePtr() *time.Time {
	return swag.Time(ValidDefaultDate())
}

// DateToString convert types.Date to string
func DateToString(v *types.Date) string {
	if v != nil {
		return v.Format(LayoutDate)
	}
	return ""
}

// DateToTime convert types.Date to time.time
func DateToTime(v *types.Date) *time.Time {
	if v != nil {
		return &v.Time
	}
	return nil
}

// TimeToString convert time.Time to string
func TimeToString(t time.Time) string {
	return t.Format(LayoutDateTime)
}

// TimeToDateString convert time.Time to string
func TimeToDateString(t time.Time) string {
	return t.Format(LayoutDate)
}

// TimeValueDefault convert *time.Time to time.time
func TimeValueDefault(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return Now()
}

// ValidFromDefaultTypesDate is default date for BiTemporal
func ValidFromDefaultTypesDate() types.Date {
	date := types.Date{}
	if err := date.UnmarshalJSON([]byte(`"0001-01-01"`)); err != nil {
		return types.Date{}
	}
	return date
}

// IsBeforeOfficialMinDateTime returns true when input value before official min value of datetime in MySQL
func IsBeforeOfficialMinDateTime(input time.Time) bool {
	officialMinValueOfDateTime := time.Date(1000, 01, 01, 0, 0, 0, 0, time.UTC)
	return input.Unix() < officialMinValueOfDateTime.Unix()
}

// PbDateToDate pb model date to time.tome
func PbDateToDate(t *date.Date) time.Time {
	return time.Date(
		int(t.GetYear()), time.Month(t.GetMonth()), int(t.GetDay()), 0, 0, 0, 0, time.UTC,
	)
}

// TimeToPbDate time to pb model date
func TimeToPbDate(t time.Time) *date.Date {
	return &date.Date{
		Year:  int32(t.Year()),
		Month: int32(t.Month()),
		Day:   int32(t.Day()),
	}
}

// Between from < t < to
func Between(t, from, to time.Time) bool {
	return from.Before(t) && to.After(t)
}

// BetweenEqual from <= t <= to
func BetweenEqual(t, from, to time.Time) bool {
	return !from.After(t) && !to.Before(t)
}

// Date struct "2006-01-02"
type Date struct {
	Format string
	time.Time
}

// UnmarshalJSON json to Date struct
func (d *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	d.Format = LayoutDate
	t, err := time.Parse(d.Format, s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// MarshalJSON Date struct to json
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(d.Format))
}

// StringDateToTime return time.Time from string date supporting for several formats
func StringDateToTime(date string) (*time.Time, error) {
	format, err := chooseFormatDate(date)
	if err != nil {
		return nil, err
	}
	dateTime, err := time.Parse(format, date)
	if err != nil {
		return nil, err
	}
	return &dateTime, nil
}

// chooseFormatDate returns format dynamically according to input date
func chooseFormatDate(date string) (string, error) {
	for _, m := range matchTypes {
		if m.regex.MatchString(date) {
			return m.format, nil
		}
	}
	return "", errors.New("not found date format")
}

var matchTypes = []struct {
	format string
	regex  *regexp.Regexp
}{
	{
		format: "2006/01/02",
		regex:  regexp.MustCompile("([0-9]{4})/(0[1-9]|1[0-2])/(0[1-9]|[12][0-9]|3[01])"),
	},
	{
		format: "2006/1/2",
		regex:  regexp.MustCompile("([0-9]{4})/([1-9])/([1-9])"),
	},
	{
		format: "2006/01/2",
		regex:  regexp.MustCompile("([0-9]{4})/(0[1-9]|1[0-2])/([1-9])"),
	},
	{
		format: "2006/1/02",
		regex:  regexp.MustCompile("([0-9]{4})/([1-9])/(0[1-9]|[12][0-9]|3[01])"),
	},
}
