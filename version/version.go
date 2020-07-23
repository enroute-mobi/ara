package version

import (
	"fmt"
	"time"
)

var value string

func Value() string {
	if value == "" {
		value = time.Now().Format("20060102-150405")
	}
	return value
}

func ApplicationName() string {
	return fmt.Sprintf("Ara %s", Value())
}
