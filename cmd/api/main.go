package main

import (
	"os"

	"github.com/ZindumKage/internal/handler"
	"github.com/ZindumKage/internal/metrics"
	"github.com/ZindumKage/internal/web"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"port": port,
		})
	})

	r.GET("/metrics", func(c *gin.Context) {
		data, err := metrics.GetMetrics()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, data)
	})

	r.POST("/jobs", handler.CreateJob)
	r.GET("/jobs/:id", handler.GetJobStatus)
	r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", web.Dashboard)
	})

	r.Run(":" + port)
}