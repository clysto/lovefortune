package main

import (
	"github.com/asdine/storm/v3"
	"github.com/clysto/lovefortune/api"
	"github.com/clysto/lovefortune/bark"
	"github.com/clysto/lovefortune/plugin"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var QweatherKey = os.Getenv("QWEATHER_KEY")

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(Logger())

	db, err := storm.Open("lovefortune.db")
	if err != nil {
		log.Fatal("数据库打开失败", err)
	}
	defer db.Close()

	taskManager := bark.NewManager(db)
	taskManager.Plugin(plugin.NewDayPlugin())
	taskManager.Plugin(plugin.NewWeatherPlugin(QweatherKey))

	var tasks []bark.BarkTask
	if err = db.All(&tasks); err != nil {
		log.Fatal("任务读取失败", err)
	}

	for index, task := range tasks {
		log.Printf("读取任务<%s@%s>", task.Description, task.ID)
		err = taskManager.AddTask(&tasks[index])
		if err != nil {
			log.Fatal("初始任务添加失败", err)
		}
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
