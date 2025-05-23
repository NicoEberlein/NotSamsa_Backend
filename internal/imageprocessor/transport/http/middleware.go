package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.New().String()
		c.Set("requestId", requestId)
		c.Set("startTime", time.Now())
		c.Next()
	}
}

func LoggingExitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetString("requestId")
		startTime, _ := c.Get("startTime")
		duration := time.Now().Sub(startTime.(time.Time))

		fmt.Printf("[%s] [RequestID: %s] - %s %s - Duration: %s\n", time.Now().UTC().Format(time.RFC3339), requestId, c.Request.URL, c.Request.Method, duration)
	}
}
