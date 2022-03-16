package api

import (
	"github.com/asdine/storm/v3"
	"github.com/clysto/lovefortune/bark"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const PageSize = 10

func ListBarkTasks(manager *bark.BarkTaskManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, manager.Tasks())
	}
}

func AddBarkTask(manager *bark.BarkTaskManager, db *storm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task bark.BarkTask
		if err := c.BindJSON(&task); err != nil {
			c.JSON(400, gin.H{
				"details": err.Error(),
			})
			return
		}
		task.ID = uuid.NewString()
		err := manager.AddTask(&task)
		if err != nil {
			c.JSON(400, gin.H{
				"details": err.Error(),
			})
			return
		}
		err = db.Save(&task)
		if err != nil {
			c.JSON(400, gin.H{
				"details": err.Error(),
			})
			return
		}
		c.JSON(200, task)
	}
}

func DeleteBarkTask(manager *bark.BarkTaskManager, db *storm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var task bark.BarkTask
		err := db.One("ID", id, &task)
		if err != nil {
			c.JSON(400, gin.H{
				"details": err.Error(),
			})
			return
		}
		err = db.DeleteStruct(&task)
		if err != nil {
			c.JSON(400, gin.H{
				"details": err.Error(),
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
		err = db.Find("TaskID", id, &logs, storm.Limit(PageSize), storm.Skip(PageSize*pageQuery.Page))
		if err != nil {
			c.JSON(400, gin.H{
				"details": err.Error(),
			})
			return
		}
		c.JSON(200, logs)
	}
}
