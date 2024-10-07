package util

import "time"

const (
	YyyyMMddhhmmss = "20060102150405"
)

func MilliFormat(millis int64, format string) string {
	return time.UnixMilli(millis).Format(format)
}
