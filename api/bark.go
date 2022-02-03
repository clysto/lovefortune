package api

import (
	"github.com/asdine/storm/v3"
	"github.com/clysto/lovefortune/bark"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const PAGE_SIZE = 10

func ListBarkTasks(manager *bark.BarkTaskManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, manager.Tasks())
	}
}

func AddBarkTask(manager *bark.BarkTaskManager, db *storm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task bark.BarkTask
		c.BindJSON(&task)
		task.ID = uuid.NewString()
		err := manager.AddTask(&task)
		if err != nil {
			c.JSON(400, gin.H{
				"description": err.Error(),
			})
			return
		}
		err = db.Save(&task)
		if err != nil {
			c.JSON(400, gin.H{
				"description": err.Error(),
			})
			return
		}
		c.JSON(200, task)
	}
}

func DeleteBarkTask(manager *bark.BarkTaskManager, db *storm.DB) gin.HandlerFunc {
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

type PageQuery struct {
	Page int
}

func GetLogs(db *storm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var pageQuery PageQuery
		err := c.BindQuery(&pageQuery)
		if err != nil {
			pageQuery.Page = 0
		} else {
			if pageQuery.Page < 1 {
				pageQuery.Page = 0
			} else {
				pageQuery.Page -= 1
			}
		}
		var logs []bark.BarkTaskLog
		err = db.Find("TaskID", id, &logs, storm.Limit(PAGE_SIZE), storm.Skip(PAGE_SIZE*pageQuery.Page))
		if err != nil {
			c.JSON(400, gin.H{
				"description": err.Error(),
			})
			return
		}
		c.JSON(200, logs)
	}
}
