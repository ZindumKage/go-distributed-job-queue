package queue

import (
	"log"
	"time"

	"github.com/ZindumKage/internal/db"
	"github.com/ZindumKage/internal/model"
)

const MaxRetries = 3

// RetryJob never blocks the caller: the re-enqueue happens after a backoff delay.
func RetryJob(jobID string) error {
	jobKey := "job:" + jobID

	// retry count lives inside the job hash (no extra keys)
	retries, err := db.Rdb.HIncrBy(db.Ctx, jobKey, "retries", 1).Result()
	if err != nil {
		return err
	}

	if retries > MaxRetries {
		return MoveToDLQ(jobID)
	}

	if err := db.Rdb.HSet(db.Ctx, jobKey, "status", "retrying").Err(); err != nil {
		return err
	}

	delay := time.Duration(retries) * time.Second
	time.AfterFunc(delay, func() {
		if err := Enqueue(model.Job{ID: jobID, Status: "retrying"}); err != nil {
			log.Println("failed to re-enqueue:", jobID, err)
		}
	})
	return nil
}