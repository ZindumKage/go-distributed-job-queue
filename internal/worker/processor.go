package worker

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
	"os"

	"github.com/ZindumKage/internal/db"
	"github.com/ZindumKage/internal/queue"
)

const WorkerPoolSize = 20

func StartWorker() {
	log.Println("Worker started, waiting for jobs...")

	jobChan := make(chan string, 100)

	for i := 0; i < WorkerPoolSize; i++ {
		go func(workerID int) {
			for jobID := range jobChan {
				log.Println("Worker", workerID, "processing:", jobID)
				process(jobID)
			}
		}(i)
	}

	for {
		// BLPop + RPush = FIFO
		res, err := db.Rdb.BLPop(db.Ctx, 0, queue.QueueName).Result()
		if err != nil {
			log.Println("Error fetching jobs:", err)
			time.Sleep(time.Second) // avoid a hot loop if Redis is down
			continue
		}
		jobChan <- res[1]
	}
}

var failRate = loadFailRate()

func loadFailRate() float64 {
	if v := os.Getenv("FAIL_RATE"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0.05
}

// shouldFail: deterministic if the job asked for it, random otherwise
func shouldFail(data map[string]string) bool {
	retries, _ := strconv.Atoi(data["retries"])
	if ft, ok := data["fail_times"]; ok {
		n, _ := strconv.Atoi(ft)
		return retries < n
	}
	return rand.Float64() < failRate
}

func process(jobID string) {
	jobKey := "job:" + jobID

	data, err := db.Rdb.HGetAll(db.Ctx, jobKey).Result()
	if err != nil || len(data) == 0 {
		log.Println("Job not found:", jobID)
		return
	}

	if err := db.Rdb.HSet(db.Ctx, jobKey, "status", "processing").Err(); err != nil {
		log.Println("Failed to update status:", err)
		return
	}

	time.Sleep(2 * time.Second) // simulate work

	if shouldFail(data) {
		attempt, _ := strconv.Atoi(data["retries"])
		log.Printf("Job %s failed (attempt %d)", jobID, attempt+1)
		db.Rdb.HSet(db.Ctx, jobKey, "last_error", fmt.Sprintf("simulated failure on attempt %d", attempt+1))
		if err := queue.RetryJob(jobID); err != nil {
			log.Println("Retry error:", jobID, err)
		}
		return
	}

	if err := db.Rdb.HSet(db.Ctx, jobKey, "status", "completed").Err(); err != nil {
		log.Println("Failed to mark completed:", err)
		return
	}
	log.Println("Job completed:", jobID)
}
