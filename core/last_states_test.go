package core

import (
	"testing"
	"time"
)

func Test_OptionParser(t *testing.T) {
	op := optionParser{}

	parsedTime := op.getOptionDuration("PT1M")
	if expected := 1 * time.Minute; parsedTime != expected {
		t.Errorf("Wrong duration parsed. Expected: %v, got: %v", expected, parsedTime)
	}

	parsedTime = op.getOptionDuration("PT30S")
	if expected := 30 * time.Second; parsedTime != expected {
		t.Errorf("Wrong duration parsed. Expected: %v, got: %v", expected, parsedTime)
	}

	parsedTime = op.getOptionDuration("PT30.000S")
	if expected := 30 * time.Second; parsedTime != expected {
		t.Errorf("Wrong duration parsed. Expected: %v, got: %v", expected, parsedTime)
	}

	parsedTime = op.getOptionDuration("PT0.1S")
	if expected := 100 * time.Millisecond; parsedTime != expected {
		t.Errorf("Wrong duration parsed. Expected: %v, got: %v", expected, parsedTime)
	}

	parsedTime = op.getOptionDuration("PT0.1S")
	if expected := 100 * time.Millisecond; parsedTime != expected {
		t.Errorf("Wrong duration parsed. Expected: %v, got: %v", expected, parsedTime)
	}

	parsedTime = op.getOptionDuration("PT0.10S")
	if expected := 100 * time.Millisecond; parsedTime != expected {
		t.Errorf("Wrong duration parsed. Expected: %v, got: %v", expected, parsedTime)
	}

	parsedTime = op.getOptionDuration("PT0.100S")
	if expected := 100 * time.Millisecond; parsedTime != expected {
		t.Errorf("Wrong duration parsed. Expected: %v, got: %v", expected, parsedTime)
	}

	parsedTime = op.getOptionDuration("PT0.001S")
	if expected := 1 * time.Millisecond; parsedTime != expected {
		t.Errorf("Wrong duration parsed. Expected: %v, got: %v", expected, parsedTime)
	}

	parsedTime = op.getOptionDuration("PT0.0001S")
	if expected := 0 * time.Minute; parsedTime != expected {
		t.Errorf("Wrong duration parsed. Expected: %v, got: %v", expected, parsedTime)
	}
}
