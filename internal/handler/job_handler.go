package handler

import (
	"net/http"
	"strconv"

	"github.com/ZindumKage/internal/db"
	"github.com/ZindumKage/internal/model"
	"github.com/ZindumKage/internal/queue"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateJob(c *gin.Context) {
	id := uuid.New().String()
	jobKey := "job:" + id

	fields := map[string]interface{}{
		"status":  "pending",
		"retries": 0,
	}

	// optional: ?fail=N makes this job fail N times before it can succeed
	if f := c.Query("fail"); f != "" {
		n, err := strconv.Atoi(f)
		if err != nil || n < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fail must be a non-negative integer"})
			return
		}
		fields["fail_times"] = n
	}

	if err := db.Rdb.HSet(db.Ctx, jobKey, fields).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save job"})
		return
	}

	if err := queue.Enqueue(model.Job{ID: id, Status: "pending"}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"job_id": id})
}

func GetJobStatus(c *gin.Context) {
	id := c.Param("id")
	jobKey := "job:" + id

	res, err := db.Rdb.HGetAll(db.Ctx, jobKey).Result()
	if err != nil || len(res) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, res)
}