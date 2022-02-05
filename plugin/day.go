package plugin

import (
	"text/template"
	"time"
)

// 2021-11-27 00:00
var WHEN_I_CONFESS = time.Date(2021, time.November, 27, 0, 0, 0, 0, time.Local)

type DayPlugin struct {
}

func (plugin *DayPlugin) LoveAnniversaryDays() int {
	duration := time.Since(WHEN_I_CONFESS)
	return int(duration.Hours() / 24)
}

func (plugin *DayPlugin) Funcs() template.FuncMap {
	return template.FuncMap{
		"loveAnniversaryDays": plugin.LoveAnniversaryDays,
	}
}

func (plugin *DayPlugin) Name() string {
	return "Day Plugin"
}

func NewDayPlugin() *DayPlugin {
	return &DayPlugin{}
}
