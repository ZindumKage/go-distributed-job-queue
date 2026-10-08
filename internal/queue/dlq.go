package queue

import "github.com/ZindumKage/internal/db"

const DLQName = "dlq"

func MoveToDLQ(jobID string) error {
	pipe := db.Rdb.TxPipeline()
	pipe.HSet(db.Ctx, "job:"+jobID, "status", "failed")
	pipe.LPush(db.Ctx, DLQName, jobID)
	_, err := pipe.Exec(db.Ctx)
	return err
}