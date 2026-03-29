package queue

import (
	"fmt"
	"time"

	"github.com/ZindumKage/internal/db"
	"github.com/ZindumKage/internal/model"
)




const MaxRetries = 3

func RetryJob(jobID string) error {
	jobKey := "job:" + jobID
	retriesKey := fmt.Sprintf("%s:retries", jobKey)

	retries, err := db.Rdb.Get(db.Ctx, retriesKey).Int()
	if err != nil {
		retries = 0
	}
	retries++

	if retries >= MaxRetries {
		return MoveToDLQ(jobID)
	}

	if err := db.Rdb.Set(db.Ctx, retriesKey, retries, 0).Err(); err != nil {
		return err
	}

	// update status to retrying
	db.Rdb.HSet(db.Ctx, jobKey, "status", "retrying")

	time.Sleep(time.Duration(retries) * time.Second)

	job := model.Job{
		ID:     jobID,
		Status: "retrying",
	}

	return Enqueue(job)
}