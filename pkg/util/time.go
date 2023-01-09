package util

import (
	"fmt"
	"strconv"
	"time"
)

const (
	Month_TIME_LAYOUT      = "2006-01"
	DEFAULT_TIME_LAYOUT    = "2006-01-02"
	DEFAULT_DB_TIME_LAYOUT = "2006-01-02 15:04:05"
)

func GetTimeLocation(location string) (*time.Location, error) {
	return time.LoadLocation(location)
}

func GetMonthsFromTimeRange(timeFrom, timeTo time.Time) ([]string, error) {
	monthInterval := SubMonth(timeTo, timeFrom)
	months := []string{}
	lastMonthFormat := timeFrom.Format(Month_TIME_LAYOUT)
	months = append(months, lastMonthFormat)
	if monthInterval <= 0 {
		return months, nil
	}
	for i := 0; i < monthInterval; i++ {
		month, err := strconv.Atoi(months[i][5:])
		if err != nil {
			return nil, err
		}
		year, err := strconv.Atoi(months[i][:4])
		if err != nil {
			return nil, err
		}
		month++
		if month > 12 {
			year++
			month = month - 12
		}
		if month < 10 {
			months = append(months, fmt.Sprintf("%d-0%d", year, month))
		} else {
			months = append(months, fmt.Sprintf("%d-%d", year, month))
		}

	}
	return months, nil
}

func SubMonth(later, earlier time.Time) (month int) {
	y1 := later.Year()
	y2 := earlier.Year()
	m1 := int(later.Month())
	m2 := int(earlier.Month())
	d1 := later.Day()
	d2 := earlier.Day()

	yearInterval := y1 - y2
	if m1 < m2 || m1 == m2 && d1 < d2 {
		yearInterval--
	}
	monthInterval := (m1 + 12) - m2
	if d1 < d2 {
		monthInterval--
	}
	monthInterval %= 12
	month = yearInterval*12 + monthInterval
	return
}
