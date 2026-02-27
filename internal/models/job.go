package models

import "time"

type JobStatus string

const (
	JobUploaded  JobStatus = "UPLOADED"
	JobRunning   JobStatus = "RUNNING"
	JobCompleted JobStatus = "COMPLETED"
	JobFailed    JobStatus = "FAILED"
)

type Job struct {
	ID             string
	Status         JobStatus
	IdempotencyKey string
	RetryCount     int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
