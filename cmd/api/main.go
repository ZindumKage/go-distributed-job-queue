package main

import (
	"os"

	"github.com/ZindumKage/internal/handler"
	"github.com/ZindumKage/internal/metrics"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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

	r.Run(":" + port)
}