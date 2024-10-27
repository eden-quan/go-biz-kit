package utils

import "time"

func NewTimeStamp() int64 {
	return time.Now().Unix()
}
