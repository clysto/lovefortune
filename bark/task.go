package bark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type BarkTaskLog struct {
	ID     string       `json:"id"`
	TaskID string       `json:"taskId" storm:"index"`
	Start  time.Time    `json:"start"`
	End    time.Time    `json:"end"`
	Buffer bytes.Buffer `json:"buffer"`
}

type BarkTask struct {
	ID          string       `json:"id"`
	ActiveID    cron.EntryID `json:"activeId"`
	Spec        string       `json:"spec"`
	Description string       `json:"description"`
	Content     string       `json:"content"`
}

func (taskLog *BarkTaskLog) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(&taskLog.Buffer, format, a...)
}

func (taskLog *BarkTaskLog) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(&taskLog.Buffer, a...)
}

func (taskLog *BarkTaskLog) UnmarshalJSON(b []byte) error {
	var log struct {
		ID     string    `json:"id"`
		TaskID string    `json:"taskId"`
		Start  time.Time `json:"start"`
		End    time.Time `json:"end"`
		Buffer string    `json:"buffer"`
	}
	err := json.Unmarshal(b, &log)
	if err != nil {
		return err
	}
	taskLog.Buffer = *bytes.NewBufferString(log.Buffer)
	taskLog.ID = log.ID
	taskLog.TaskID = log.TaskID
	taskLog.Start = log.Start
	taskLog.End = log.End
	return nil
}

func (taskLog *BarkTaskLog) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID     string    `json:"id"`
		TaskID string    `json:"taskId"`
		Start  time.Time `json:"start"`
		End    time.Time `json:"end"`
		Buffer string    `json:"buffer"`
	}{
		ID:     taskLog.ID,
		TaskID: taskLog.TaskID,
		Start:  taskLog.Start,
		End:    taskLog.End,
		Buffer: taskLog.Buffer.String(),
	})
}
