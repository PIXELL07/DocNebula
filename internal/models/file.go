package models

import "time"

type FileStatus string

const (
	FilePending    FileStatus = "PENDING"
	FileUnzipped   FileStatus = "UNZIPPED"
	FileOCRDone    FileStatus = "OCR_DONE"
	FileVectorDone FileStatus = "VECTOR_DONE"
	FileFailed     FileStatus = "FAILED"
)

// File represents a single document extracted from the uploaded ZIP
// This is the unit that flows through the OCR → Vector pipeline.
type File struct {
	ID        string
	JobID     string
	Path      string // object storage path
	Status    FileStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
