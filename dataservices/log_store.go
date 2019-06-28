package dataservices

import (
	"github.com/bitrise-io/addons-ship-backend/models"
)

// LogStore ...
type LogStore interface {
	Get(string) (models.LogChunk, error)
	Set(string, models.LogChunk) error
}
