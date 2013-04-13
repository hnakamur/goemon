package time

import (
	"errors"
	"fmt"
	"regexp"
)
import T "time"

const (
  DAYS_IN_MONTH = 31
  DAYS_IN_YEAR = 366
  HOURS_IN_DAY = 24
  DAY_TEXT = "d"
  MONTH_TEXT = "M"
  YEAR_TEXT = "Y"
)

var re *regexp.Regexp

func init() {
	re = regexp.MustCompile("(\\d+(?:\\.\\d*)?|\\.\\d+)(Y|M|d|h|m|s|us|ns)")
}

func ParseDuration(s string) (T.Duration, error) {
	matches := re.FindAllStringSubmatch(s, -1)
	if matches == nil {
		return 0, errors.New("invalid duration string")
	}
	var duration float64
	totalLen := 0
	var val float64
	for _, match := range matches {
		totalLen += len(match[0])

		n, err := fmt.Sscanf(match[1], "%f", &val)
		if err != nil {
			return 0, err
		}
		if n != 1 {
			return 0, errors.New("multiple match error")
		}

		unit, err := unitToDuration(match[2])
		if err != nil {
			return 0, err
		}

		duration += val * float64(unit)
	}
	if totalLen != len(s) {
		return 0, errors.New("invalid duration string")
	}
	return T.Duration(duration), nil
}

func unitToDuration(s string) (T.Duration, error) {
	switch s {
	case "Y": return DAYS_IN_YEAR * HOURS_IN_DAY * T.Hour, nil
	case "M": return DAYS_IN_MONTH * HOURS_IN_DAY * T.Hour, nil
	case "d": return HOURS_IN_DAY * T.Hour, nil
	case "h": return T.Hour, nil
	case "m": return T.Minute, nil
	case "s": return T.Second, nil
	case "ms": return T.Millisecond, nil
	case "us": return T.Microsecond, nil
	case "ns": return T.Nanosecond, nil
	}
	return 0, errors.New("Invalid duration unit")
}
