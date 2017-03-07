package model

import (
	"fmt"
	"time"
)

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func NewDate(reference time.Time) Date {
	return Date{
		Year:  reference.Year(),
		Month: reference.Month(),
		Day:   reference.Day(),
	}
}

func (date *Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", date.Year, date.Month, date.Day)
}
