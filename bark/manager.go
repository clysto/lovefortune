package bark

import (
	"bytes"
	"text/template"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/clysto/lovefortune/plugin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type BarkTaskManager struct {
	*cron.Cron
	tasks []*BarkTask
	db    *storm.DB
}

func NewManager(db *storm.DB) *BarkTaskManager {
	c := cron.New()
	manager := &BarkTaskManager{
		Cron:  c,
		tasks: nil,
		db:    db,
	}
	return manager
}

func taskFunc(task *BarkTask, taskLog *BarkTaskLog) {
	var sendContent string
	taskLog.Println("开始执行: " + task.Description)
	taskLog.Println("发送内容")
	taskLog.Println("=====================")
	t, err := template.New("tmp").Parse(task.Content)
	if err != nil {
		sendContent = task.Content
	} else {
		var buf bytes.Buffer
		err = t.Execute(&buf, gin.H{
			"LoveAnniversaryDays": plugin.LoveAnniversaryDays(),
		})
		if err != nil {
			sendContent = task.Content
		} else {
			sendContent = buf.String()
		}
	}
	taskLog.Println(sendContent)
	taskLog.Println("=====================")
}

func (manager *BarkTaskManager) AddTask(task *BarkTask) (cron.EntryID, error) {
	id, err := manager.AddFunc(task.Spec, func() {
		var taskLog BarkTaskLog
		taskLog.ID = uuid.NewString()
		taskLog.TaskID = task.ID
		taskLog.Start = time.Now()
		taskFunc(task, &taskLog)
		taskLog.End = time.Now()
		manager.db.Save(&taskLog)
	})
	if err == nil {
		task.ActiveID = id
		manager.tasks = append(manager.tasks, task)
	}
	return id, err
}

func (manager *BarkTaskManager) Tasks() []*BarkTask {
	return manager.tasks
}

func (manager *BarkTaskManager) RemoveTask(id string) {
	var tasks []*BarkTask
	var activeID cron.EntryID
	for _, task := range manager.tasks {
		if task.ID != id {
			tasks = append(tasks, task)
		} else {
			activeID = task.ActiveID
		}
	}
	manager.Remove(activeID)
	manager.tasks = tasks
}
