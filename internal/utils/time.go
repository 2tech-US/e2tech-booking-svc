package utils

import "time"

var layoutDate = "2006/01/02"

func ParseStringToDate(date string) (time.Time, error) {
	return time.Parse(layoutDate, date)
}

func ParsedDateToString(date time.Time) string {
	return date.Format(layoutDate)
}
