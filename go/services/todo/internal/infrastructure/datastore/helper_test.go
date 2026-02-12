package datastore_test

import "time"

func getLocalTimeByString(expectedDateStr string) time.Time {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	layout := "2006-01-02T15:04:00Z"
	expectedDatetime, _ := time.ParseInLocation(layout, expectedDateStr, loc)
	return expectedDatetime
}
