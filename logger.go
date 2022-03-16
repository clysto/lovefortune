package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

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
