package queue

import "time"

// PageMessage represents page-level OCR work.
type PageMessage struct {
	JobID     string    `json:"job_id"`
	FileID    string    `json:"file_id"`
	PageNum   int       `json:"page_num"`
	ImagePath string    `json:"image_path"`
	Attempt   int       `json:"attempt"`
	Timestamp time.Time `json:"ts"`
}
