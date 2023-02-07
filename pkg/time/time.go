package pkg

import "time"

func ToTimeFromUnix(unix int64) time.Time {
	return time.Unix(unix, 0)
}

func ToUnixTime(t time.Time) int64 {
	return t.Unix()
}
