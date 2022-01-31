package main

import (
	"log"

	"github.com/asdine/storm/v3"
	"github.com/clysto/lovefortune/bark"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func listBarkTasks(manager *bark.BarkTaskManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, manager.Tasks())
	}
}

func addBarkTask(manager *bark.BarkTaskManager, db *storm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task bark.BarkTask
		c.BindJSON(&task)
		task.ID = uuid.NewString()
		err := db.Save(&task)
		if err != nil {
			c.JSON(400, gin.H{
				"description": err.Error(),
			})
			return
		}
		manager.AddTask(&task)
		c.JSON(200, task)
	}
}

func deleteBarkTask(manager *bark.BarkTaskManager, db *storm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		print(id)
		var task bark.BarkTask
		err := db.One("ID", id, &task)
		if err != nil {
			c.JSON(400, gin.H{
				"description": err.Error(),
			})
			return
		}
		err = db.DeleteStruct(&task)
		if err != nil {
			c.JSON(400, gin.H{
				"description": err.Error(),
			})
			return
		}
		manager.RemoveTask(task.ID)
		c.JSON(204, task)
	}
}

func getLogs(db *storm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var logs []bark.BarkTaskLog
		err := db.Find("TaskID", id, &logs)
		if err != nil {
			c.JSON(400, gin.H{
				"description": err.Error(),
			})
			return
		}
		c.JSON(200, logs)
	}
}

func main() {
	router := gin.Default()

	db, err := storm.Open("lovefortune.db")
	if err != nil {
		log.Fatal("数据库打开失败")
	}
	defer db.Close()

	taskManager := bark.NewManager(db)

	var tasks []bark.BarkTask
	err = db.All(&tasks)
	if err != nil {
		log.Fatal("任务读取失败")
	}

	for _, task := range tasks {
		taskManager.AddTask(&task)
	}
	taskManager.Start()
	defer taskManager.Stop()

	api := router.Group("/api")
	api.GET("/tasks", listBarkTasks(taskManager))
	api.POST("/tasks", addBarkTask(taskManager, db))
	api.DELETE("/tasks/:id", deleteBarkTask(taskManager, db))
	api.GET("/logs/:id", getLogs(db))
	router.Run()
}
