package repo

import (
	"encoding/base64"
	"strconv"
	"time"
)

// DecodeCursor will decode cursor from user for db
func DecodeCursor(encodedTime string) (time.Time, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Unix(1<<31-1, 0), err
	} else if len(byt) == 0 {
		return time.Unix(1<<31-1, 0), nil
	}

	timeString := string(byt)
	millis, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		return time.Unix(1<<31-1, 0), err
	}

	return time.Unix(0, millis*int64(time.Millisecond)), nil
}

// EncodeCursor will encode cursor from db to user
func EncodeCursor(t time.Time) string {
	timeString := strconv.FormatInt(t.UnixMilli(), 10)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}
