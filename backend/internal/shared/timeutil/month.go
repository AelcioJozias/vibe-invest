package timeutil

import (
	"fmt"
	"time"
)

const referenceMonthLayout = "2006-01"

func ParseReferenceMonth(value string) (time.Time, error) {
	month, err := time.Parse(referenceMonthLayout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("referenceMonth must be in YYYY-MM format")
	}

	return ToMonthDate(month), nil
}

func CurrentReferenceMonth(now time.Time) time.Time {
	return ToMonthDate(now)
}

func PreviousMonth(month time.Time) time.Time {
	m := ToMonthDate(month)
	return m.AddDate(0, -1, 0)
}

func ToMonthDate(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, time.UTC)
}
