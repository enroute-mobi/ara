package ls

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type optionParser struct{}

func (parser *optionParser) getOptionDuration(option string) time.Duration {

	durationRegex := regexp.MustCompile(`P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)(?:[\.,](\d{1,3}))?S)?)?`)
	matches := durationRegex.FindStringSubmatch(strings.TrimSpace(option))

	if len(matches) == 0 {
		return 0
	}
	years := parser.parseDuration(matches[1]) * 24 * 365 * time.Hour
	months := parser.parseDuration(matches[2]) * 30 * 24 * time.Hour
	days := parser.parseDuration(matches[3]) * 24 * time.Hour
	hours := parser.parseDuration(matches[4]) * time.Hour
	minutes := parser.parseDuration(matches[5]) * time.Minute
	seconds := parser.parseDuration(matches[6]) * time.Second
	rest := parser.parseDuration(matches[7]) * parser.durationPow(10, 3-len(matches[7])) * time.Millisecond

	return time.Duration(years + months + days + hours + minutes + seconds + rest)
}

func (parser *optionParser) parseDuration(value string) time.Duration {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return time.Duration(parsed)
}

func (parser *optionParser) durationPow(a, b int) time.Duration {
	p := 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return time.Duration(p)
}
