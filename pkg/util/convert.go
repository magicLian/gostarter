package util

import (
	"strconv"
	"time"
)

func SetDefaultFloat64(targetstr string, defaultVal ...float64) float64 {
	val, err := strconv.ParseFloat(targetstr, 64)
	if len(defaultVal) > 0 && err != nil {
		targetstr = strconv.FormatFloat(defaultVal[0], 'f', -1, 64)
		return defaultVal[0]
	}
	return val
}

func SetDefaultDuration(targetstr string, defaultVal ...time.Duration) time.Duration {
	val, err := time.ParseDuration(targetstr)
	if len(defaultVal) > 0 && err != nil {
		targetstr = defaultVal[0].String()
		return defaultVal[0]
	}
	return val
}

func SetDefaultString(targetstr string, defaultVal string) string {
	if targetstr == "" {
		return defaultVal
	}
	return targetstr
}

func SetDefaultBool(targetstr string, defaultVal bool) bool {
	if targetstr == "true" {
		return true
	} else if targetstr == "false" {
		return false
	}
	return defaultVal
}

func SetDefaultInt(targetstr string, defaultVal int) int {
	number, err := strconv.Atoi(targetstr)
	if targetstr != "" && err == nil {
		return number
	}
	return defaultVal

}

func SetDefaultInt64(targetstr string, defaultVal int64) int64 {
	number, err := strconv.ParseInt(targetstr, 10, 64)
	if targetstr != "" && err == nil {
		return number
	}
	return defaultVal

}
