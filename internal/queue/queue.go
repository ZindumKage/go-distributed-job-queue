package queue

import (
	"github.com/ZindumKage/internal/db"
	"github.com/ZindumKage/internal/model"
)

const QueueName = "jobs"

func Enqueue(job model.Job) error {
	return db.Rdb.RPush(db.Ctx, QueueName, job.ID).Err()
}