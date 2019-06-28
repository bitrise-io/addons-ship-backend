package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// LogChunk ...
type LogChunk struct {
	ID        uuid.UUID `json:"id"`
	TaskID    uuid.UUID `json:"task_id"`
	Pos       int       `json:"pos"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
