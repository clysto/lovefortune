package main

import (
	"log"

	"github.com/asdine/storm/v3"
	"github.com/clysto/lovefortune/api"
	"github.com/clysto/lovefortune/bark"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 允许跨域
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("X-Access-Token")
	router.Use(cors.New(config))

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

	apiRouter := router.Group("/api")

	apiRouter.Use(api.Auth())
	apiRouter.GET("/tasks", api.ListBarkTasks(taskManager))
	apiRouter.POST("/tasks", api.AddBarkTask(taskManager, db))
	apiRouter.DELETE("/tasks/:id", api.DeleteBarkTask(taskManager, db))
	apiRouter.GET("/logs/:id", api.GetLogs(db))

	router.Run()
}
