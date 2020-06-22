package logparser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var timeChars = map[string]string{
	"m": "minute(s)",
	"s": "second(s)",
	"h": "hour(s)",
}

var durationRe = regexp.MustCompile(`\d+`)

func ParseDuration(start, s string) (time.Time, error) {

	result := durationRe.FindSubmatch([]byte(s))
	unit := string(s[len(s)-1:])

	timeValue, err := strconv.Atoi(string(result[0]))

	fmt.Printf("We are adding %v %v to the start time...", timeValue, timeChars[unit])

	value, err := time.Parse("2006-01-02T15:04:05", start)

	switch unit {
	case "h":
		value = value.Add(time.Duration(timeValue) * time.Hour)
	case "m":
		value = value.Add(time.Duration(timeValue) * time.Minute)
	case "s":
		value = value.Add(time.Duration(timeValue) * time.Second)
	}

	return value, err
}
