package valueObject

import "time"

type UnixTime int64

func (ut UnixTime) GetUnixTime() int64 {
	return time.Unix(int64(ut), 0).UTC().Unix()
}

func (ut UnixTime) GetRfcDate() string {
	return time.Unix(int64(ut), 0).UTC().Format(time.RFC3339)
}

func (ut UnixTime) GetDateOnly() string {
	return time.Unix(int64(ut), 0).UTC().Format("2006-01-02")
}

func (ut UnixTime) GetTimeOnly() string {
	return time.Unix(int64(ut), 0).UTC().Format("15:04:05")
}
