package worker

import (
	"log"
	"time"

	"github.com/ZindumKage/internal/db"
	"github.com/ZindumKage/internal/queue"
)




const QueueName = "jobs"

func StartWorker() {
	log.Println("Worker started, waiting for jobs...")

	for {
		res, err := db.Rdb.BRPop(db.Ctx, 0*time.Second, QueueName).Result()
		if err != nil {
			log.Println("Error fetching jobs:", err)
			continue
		}
		jobID := res[1]

		go process(jobID)
	}
	}

	func process(jobID string) {
	log.Println("Processing job:", jobID)

	jobKey := "job:" + jobID

	
if err := db.Rdb.HSet(db.Ctx, jobKey, "status", "processing").Err(); err != nil {
	log.Println("Failed to update status:", err)
	return
}

	time.Sleep(2 * time.Second)

	
	if len(jobID) > 0 && jobID[len(jobID)-1] == '0' {
		log.Println("Job failed:", jobID)
		queue.RetryJob(jobID)
		return
	}

	
	db.Rdb.HSet(db.Ctx, jobKey, "status", "completed")
	log.Println("Job completed:", jobID)
}