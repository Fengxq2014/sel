package router

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)
		log.Println(latency)
		status := c.Writer.Status()
		log.Println(status)
	}
}
