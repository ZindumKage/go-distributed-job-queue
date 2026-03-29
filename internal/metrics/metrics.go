package metrics

import (
	"github.com/ZindumKage/internal/db"
)

type Metrics struct {
	TotalJobs      int `json:"total_jobs"`
	PendingJobs    int `json:"pending_jobs"`
	ProcessingJobs int `json:"processing_jobs"`
	CompletedJobs  int `json:"completed_jobs"`
	FailedJobs     int `json:"failed_jobs"`
	DLQSize        int `json:"dlq_size"`
}

func GetMetrics() (*Metrics, error) {
	keys, err := db.Rdb.Keys(db.Ctx, "job:*").Result()
	if err != nil {
		return nil, err
	}

	var m Metrics

	for _, key := range keys {
		status, err := db.Rdb.HGet(db.Ctx, key, "status").Result()
		if err != nil {
			continue
		}
		switch status {
		case "pending":
			m.PendingJobs++
		case "processing":
			m.ProcessingJobs++
		case "completed":
			m.CompletedJobs++
		case "failed":
			m.FailedJobs++

		}

		m.TotalJobs++
	}

	dlqSize, _ := db.Rdb.LLen(db.Ctx, "dlq").Result()

	m.DLQSize = int(dlqSize)

	return &m, nil
}
