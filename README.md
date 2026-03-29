# Distributed Job Queue Service (Go)

## Overview

This project is a distributed job processing system built in Go. It demonstrates how to design and scale a backend system using a queue-based architecture with horizontal scaling, load balancing, and fault tolerance.

The system accepts jobs via an API, queues them in Redis, processes them asynchronously using worker pools, and exposes metrics for monitoring.

---

## Architecture

Client → Nginx Load Balancer → Multiple API Instances → Redis Queue → Worker Processes → Job Results

---

## Features

- REST API for job submission and status tracking  
- Redis-backed job queue  
- Worker pool with concurrent processing  
- Horizontal scaling for API and workers  
- Load balancing using Nginx (round-robin)  
- Retry mechanism for failed jobs  
- Dead Letter Queue (DLQ) for persistent failures  
- Metrics endpoint for monitoring system state  

---

## Tech Stack

- Go (Golang)
- Gin (HTTP framework)
- Redis (queue and storage)
- Nginx (reverse proxy and load balancer)

---

## Project Structure

job-service/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── queue/
│   ├── worker/
│   ├── metrics/
│   ├── db/
│   └── model/
├── pkg/
│   └── utils/
├── nginx/
│   └── nginx.conf
├── Dockerfile (optional)
└── docker-compose.yml (optional)

---

## Getting Started

### Prerequisites

- Go installed
- Redis running locally
- Nginx installed

---

## Running the System

### 1. Start Redis

redis-server

---

### 2. Start API Instances (Horizontal Scaling)

PORT=8080 go run cmd/api/main.go &
PORT=8081 go run cmd/api/main.go &
PORT=8082 go run cmd/api/main.go &

---

### 3. Start Worker Processes

for i in {1..5}; do go run cmd/worker/main.go & done

Each worker process runs a pool of concurrent workers.

---

### 4. Configure Nginx

events {}

http {
    upstream backend {
        server 127.0.0.1:8080;
        server 127.0.0.1:8081;
        server 127.0.0.1:8082;
    }

    server {
        listen 8000;

        location / {
            proxy_pass http://backend;
        }
    }
}

Reload Nginx:

nginx -t
nginx -s reload

---

## API Endpoints

### Health Check

GET /health

Response:

{
  "status": "ok",
  "port": "8080"
}

---

### Create Job

POST /jobs

---

### Get Job Status

GET /jobs/:id

---

### Metrics

GET /metrics

Example response:

{
  "total_jobs": 1000,
  "pending_jobs": 0,
  "processing_jobs": 0,
  "completed_jobs": 950,
  "failed_jobs": 50,
  "dlq_size": 50
}

---

## Load Testing

### Sequential

time for i in {1..1000}; do
  curl -s -X POST http://localhost:8000/jobs > /dev/null
done

### Parallel (Recommended)

seq 1000 | xargs -P 20 -I {} curl -s -X POST http://localhost:8000/jobs > /dev/null

---

## Clearing Redis (for clean tests)

redis-cli FLUSHALL

---

## Key Concepts Demonstrated

- Horizontal scaling (multiple API instances)
- Load balancing (Nginx round-robin)
- Asynchronous job processing
- Worker pool concurrency
- Retry logic and failure handling
- Dead Letter Queue (DLQ)
- Distributed system design

---

## Limitations

- No authentication or authorization
- No persistent storage beyond Redis
- No auto-scaling

---

## Future Improvements

- Add rate limiting  
- Introduce monitoring (Prometheus, Grafana)  
- Containerize with Docker  
- Deploy with Kubernetes  
- Add health checks for worker nodes  

---

## License

MIT License