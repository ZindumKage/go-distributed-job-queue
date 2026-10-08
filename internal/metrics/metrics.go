package metrics

import (
	"github.com/ZindumKage/internal/db"
	"github.com/redis/go-redis/v9"
)

type Metrics struct {
	TotalJobs      int `json:"total_jobs"`
	PendingJobs    int `json:"pending_jobs"`
	ProcessingJobs int `json:"processing_jobs"`
	RetryingJobs   int `json:"retrying_jobs"`
	CompletedJobs  int `json:"completed_jobs"`
	FailedJobs     int `json:"failed_jobs"`
	DLQSize        int `json:"dlq_size"`
}

func GetMetrics() (*Metrics, error) {
	var m Metrics
	seen := map[string]struct{}{} // SCAN may return duplicates
	var cursor uint64

	for {
		keys, next, err := db.Rdb.Scan(db.Ctx, cursor, "job:*", 500).Result()
		if err != nil {
			return nil, err
		}

		var batch []string
		for _, k := range keys {
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				batch = append(batch, k)
			}
		}

		pipe := db.Rdb.Pipeline()
		cmds := make([]*redis.StringCmd, len(batch))
		for i, k := range batch {
			cmds[i] = pipe.HGet(db.Ctx, k, "status")
		}
		pipe.Exec(db.Ctx) // individual errors are checked per command below

		for _, c := range cmds {
			status, err := c.Result()
			if err != nil {
				continue
			}
			m.TotalJobs++
			switch status {
			case "pending":
				m.PendingJobs++
			case "processing":
				m.ProcessingJobs++
			case "retrying":
				m.RetryingJobs++
			case "completed":
				m.CompletedJobs++
			case "failed":
				m.FailedJobs++
			}
		}

		cursor = next
		if cursor == 0 {
			break
		}
	}

	dlq, _ := db.Rdb.LLen(db.Ctx, "dlq").Result()
	m.DLQSize = int(dlq)
	return &m, nil
}