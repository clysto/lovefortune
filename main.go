package main

import (
	"log"
	"os"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/clysto/lovefortune/api"
	"github.com/clysto/lovefortune/bark"
	"github.com/clysto/lovefortune/plugin"
	"github.com/gin-gonic/gin"
)

var QWEATHER_KEY = os.Getenv("QWEATHER_KEY")

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		TimeStamp := time.Now()
		Latency := TimeStamp.Sub(start)

		ClientIP := c.ClientIP()
		Method := c.Request.Method
		StatusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		Path := path

		log.Printf("%3d | %13v | %15s | %-7s  %#v\n",
			StatusCode,
			Latency,
			ClientIP,
			Method,
			Path,
		)
	}
}

func main() {
	router := gin.New()
	gin.DisableConsoleColor()
	router.Use(gin.Recovery())
	router.Use(Logger())

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

	for index, task := range tasks {
		log.Printf("读取任务<%s@%s>", task.Description, task.ID)
		err = taskManager.AddTask(&tasks[index])
		if err != nil {
			log.Fatal("初始任务添加失败")
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
