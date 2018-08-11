package helper

import "time"

func GetUnixTimestamp() (int64) {
	currentTime := time.Now()
	unixTime := currentTime.Unix()

	return unixTime
}
