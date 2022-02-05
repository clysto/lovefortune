package main

import (
	"log"
	"os"

	"github.com/asdine/storm/v3"
	"github.com/clysto/lovefortune/api"
	"github.com/clysto/lovefortune/bark"
	"github.com/clysto/lovefortune/plugin"
	"github.com/gin-gonic/gin"
)

var QWEATHER_KEY = os.Getenv("QWEATHER_KEY")

func main() {
	router := gin.New()
	gin.DisableConsoleColor()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	db, err := storm.Open("lovefortune.db")
	if err != nil {
		log.Fatal("数据库打开失败")
	}
	defer db.Close()

	taskManager := bark.NewManager(db)
	taskManager.Plugin(plugin.NewDayPlugin())
	taskManager.Plugin(plugin.NewWeatherPlugin(QWEATHER_KEY))

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
