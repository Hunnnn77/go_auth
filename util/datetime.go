package util

import "time"

var loc, _ = time.LoadLocation("Asia/Tokyo")
var NowInTimezone = time.Now().In(loc)
