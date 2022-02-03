package bark

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
	c := cron.New(cron.WithSeconds())
	manager := &BarkTaskManager{
		Cron:  c,
		tasks: nil,
		db:    db,
	}
	return manager
}

func taskFunc(task *BarkTask, taskLog *BarkTaskLog) {
	var sendContent string
	taskLog.Println("开始执行:", task.Description)
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
	taskLog.Println("发送标题:", task.Title)
	taskLog.Println("发送内容:", sendContent)
	taskLog.Println("显示图标:", task.Icon)
	body := BarkPushBody{
		Body:      sendContent,
		Title:     task.Title,
		DeviceKey: task.DeviceKey,
		ExtParams: BarkPushBodyExtParams{
			Icon: task.Icon,
		},
	}
	json, err := json.Marshal(body)
	if err != nil {
		taskLog.Println(err)
		return
	}
	buf := bytes.NewBuffer(json)
	client := &http.Client{}

	// 创建 http client
	req, err := http.NewRequest(http.MethodPost, "https://api.day.app/push", buf)
	if err != nil {
		taskLog.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// 发送请求
	resp, err := client.Do(req)

	if err != nil {
		taskLog.Println(err)
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		taskLog.Println(err)
		return
	}

	taskLog.Println("请求结果:")
	taskLog.Println("Response Status    :", resp.Status)
	taskLog.Println("Response Headers   :", resp.Header)
	taskLog.Println("Response Body      :", string(respBody))
	taskLog.Println("结束执行:", task.Description)
}

func (manager *BarkTaskManager) AddTask(task *BarkTask) error {
	id, err := manager.AddFunc(task.Spec, func() {
		var taskLog BarkTaskLog
		taskLog.ID = uuid.NewString()
		taskLog.TaskID = task.ID
		taskLog.Start = time.Now()
		// 执行发送任务
		taskFunc(task, &taskLog)
		taskLog.End = time.Now()
		manager.db.Save(&taskLog)
	})
	if err == nil {
		task.ActiveID = id
		manager.tasks = append(manager.tasks, task)
	}
	return err
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
