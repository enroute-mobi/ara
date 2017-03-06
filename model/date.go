package model

import "time"

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
