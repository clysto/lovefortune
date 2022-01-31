package plugin

import (
	"time"
)

// 2021-11-27 00:00
var WHEN_I_CONFESS = time.Date(2021, time.November, 27, 0, 0, 0, 0, time.Local)

func LoveAnniversaryDays() int {
	duration := time.Since(WHEN_I_CONFESS)
	return int(duration.Hours() / 24)
}
