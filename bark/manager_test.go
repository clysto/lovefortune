package bark

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/clysto/lovefortune/plugin"
	"github.com/google/uuid"
)

var (
	QWEATHER_KEY = os.Getenv("QWEATHER_KEY")
	DEVICE_KEY   = os.Getenv("DEVICE_KEY")
)

func TestDay(t *testing.T) {
	taskManager := NewManager(nil)
	taskManager.Plugin(plugin.NewDayPlugin())
	taskManager.Plugin(plugin.NewWeatherPlugin(QWEATHER_KEY))
	task := BarkTask{
		ID:          uuid.NewString(),
		ActiveID:    1,
		Spec:        "0 0 0 * * *",
		Description: "test",
		Title:       "test",
		Content:     "在一起第{{ loveAnniversaryDays }}天",
		DeviceKey:   DEVICE_KEY,
	}
	var taskLog BarkTaskLog
	taskLog.ID = uuid.NewString()
	taskLog.TaskID = task.ID
	taskLog.Start = time.Now()
	taskManager.TaskFunc(&task, &taskLog)
	taskLog.End = time.Now()
	fmt.Print(taskLog.Buffer.String())
}

func TestWeather(t *testing.T) {
	taskManager := NewManager(nil)
	taskManager.Plugin(plugin.NewDayPlugin())
	taskManager.Plugin(plugin.NewWeatherPlugin(QWEATHER_KEY))
	content := "{{ $w := (weather \"101220601\").Now }}" +
		"现在天气：{{$w.Temp}}度，体感温度{{$w.FeelsLike}}度；" +
		"{{$w.Text}}，{{$w.WindDir}}{{$w.WindScale}}级；" +
		"湿度：{{$w.Humidity}}%"
	task := BarkTask{
		ID:          uuid.NewString(),
		ActiveID:    1,
		Spec:        "0 0 0 * * *",
		Description: "test",
		Title:       "test",
		Content:     content,
		DeviceKey:   DEVICE_KEY,
	}
	var taskLog BarkTaskLog
	taskLog.ID = uuid.NewString()
	taskLog.TaskID = task.ID
	taskLog.Start = time.Now()
	taskManager.TaskFunc(&task, &taskLog)
	taskLog.End = time.Now()
	fmt.Print(taskLog.Buffer.String())
}
