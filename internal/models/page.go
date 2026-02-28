package models

import "time"

// per‑page OCR progress so workers can resume safely.
type Page struct {
	ID        string
	FileID    string
	PageNum   int
	Text      string
	Done      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
