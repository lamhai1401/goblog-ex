package utils

import "time"

// NowStr parsing now to string with format
func NowStr() string {
	return time.Now().Format("2006-01-02T15:04:05.000")
}
