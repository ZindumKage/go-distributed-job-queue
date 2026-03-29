package queue

import (
	"encoding/json"

	"github.com/ZindumKage/internal/db"
	"github.com/ZindumKage/internal/model"
)




const QueueName = "jobs"

func Enqueue(job model.Job) error {
	data, _ := json.Marshal(job)
	return db.Rdb.RPush(db.Ctx, QueueName, data).Err()
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