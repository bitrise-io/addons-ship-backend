package analytics

import (
	"errors"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	segment "gopkg.in/segmentio/analytics-go.v3"
)

// Interface ...
type Interface interface {
	FirstVersionCreated(appSlug, buildSlug, platform string)
	PublishFinished(appSlug string, appVersionID uuid.UUID, result string)
}

// Client ...
type Client struct {
	client segment.Client
	logger *zap.Logger
}

// NewClient ...
func NewClient(logger *zap.Logger) (Client, error) {
	writeKey, ok := os.LookupEnv("SEGMENT_WRITE_KEY")
	if !ok {
		return Client{}, errors.New("No value set for env SEGMENT_WRITEKEY")
	}

	return Client{
		client: segment.New(writeKey),
		logger: logger,
	}, nil
}

// FirstVersionCreated ...
func (c *Client) FirstVersionCreated(appSlug, buildSlug, platform string) {
	err := c.client.Enqueue(segment.Track{
		UserId: appSlug,
		Event:  "First app version was created",
		Properties: segment.NewProperties().
			Set("app_slug", appSlug).
			Set("build_slug", buildSlug).
			Set("platform", platform).
			Set("datetime", time.Now()),
	})
	if err != nil {
		c.logger.Warn("Failed to track analytics (FirstVersionCreated)", zap.Error(err))
	}
}

// PublishFinished ...
func (c *Client) PublishFinished(appSlug string, appVersionID uuid.UUID, result string) {
	err := c.client.Enqueue(segment.Track{
		UserId: appSlug,
		Event:  "Publish task finished",
		Properties: segment.NewProperties().
			Set("app_slug", appSlug).
			Set("app_version_id", appVersionID).
			Set("result", result).
			Set("datetime", time.Now()),
	})
	if err != nil {
		c.logger.Warn("Failed to track analytics (PublishFinished)", zap.Error(err))
	}
}
