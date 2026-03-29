package worker

import (
	"log"
	"time"

	"github.com/ZindumKage/internal/db"
	"github.com/ZindumKage/internal/queue"
)

const QueueName = "jobs"
const WorkerPoolSize = 20

func StartWorker() {
	log.Println("Worker started, waiting for jobs...")

	jobChan := make(chan string, 100)

	//  Start worker pool ONCE
	for i := 0; i < WorkerPoolSize; i++ {
		go func(workerID int) {
			for jobID := range jobChan {
				log.Println("Worker", workerID, "processing:", jobID)
				process(jobID)
			}
		}(i)
	}

	// Listen for jobs
	for {
		res, err := db.Rdb.BRPop(db.Ctx, 0*time.Second, QueueName).Result()
		if err != nil {
			log.Println("Error fetching jobs:", err)
			continue
		}

		jobID := res[1]

		// send to worker pool
		jobChan <- jobID
	}
}

func process(jobID string) {
	log.Println("Processing job:", jobID)

	jobKey := "job:" + jobID

	// fetch job data
	jobData, err := db.Rdb.HGetAll(db.Ctx, jobKey).Result()
	if err != nil || len(jobData) == 0 {
		log.Println("Job not found:", jobID)
		return
	}

	// update status → processing
	if err := db.Rdb.HSet(db.Ctx, jobKey, "status", "processing").Err(); err != nil {
		log.Println("Failed to update status:", err)
		return
	}

	time.Sleep(2 * time.Second)

	// simulate failure
	if len(jobID) > 0 && jobID[len(jobID)-1] == '0' {
		log.Println("Job failed:", jobID)
		time.Sleep(2 * time.Second)
		queue.RetryJob(jobID)
		return
	}

	// success
	if err := db.Rdb.HSet(db.Ctx, jobKey, "status", "completed").Err(); err != nil {
		log.Println("Failed to mark completed:", err)
		return
	}

	log.Println("Job completed:", jobID)
}