package queue

import "github.com/ZindumKage/internal/db"


const DLQName = "dlq"

func MoveToDLQ(jobID string) error {
	jobKey := "job:" + jobID

	db.Rdb.HSet(db.Ctx, jobKey, "status", "failed")

	return db.Rdb.LPush(db.Ctx, DLQName, jobID).Err()
}