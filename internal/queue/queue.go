package queue

import (
	"encoding/json"

	"github.com/ZindumKage/internal/db"
	"github.com/ZindumKage/internal/model"
)




const QueueName = "jobs"

func Enqueue(job model.Job) error {
	
	return db.Rdb.RPush(db.Ctx, QueueName, job.ID).Err()
}

func Dequeue() (*model.Job, error) {
	res, err := db.Rdb.BLPop(db.Ctx, 0, QueueName).Result()
	if err != nil {
		return nil, err
	}
	var job model.Job
	json.Unmarshal([]byte(res[1]), &job)

	return &job, nil
}