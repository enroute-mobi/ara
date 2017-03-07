package model

import "testing"

func Test_Date_String(t *testing.T) {
	date := Date{Year: 1, Month: 2, Day: 3}
	if expected := "0001-02-03"; date.String() != expected {
		t.Errorf("Date.String() returns wrong value, got: %s, required: %s", date.String(), expected)
	}
}
