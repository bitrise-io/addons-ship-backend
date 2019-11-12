package services_test

import uuid "github.com/satori/go.uuid"

type testAnalyticsClient struct {
	firstVersionCreatedFn func(appSlug, buildSlug, platform string)
	publishFinishedFn     func(appSlug string, appVersionID uuid.UUID, result string)
}

func (c *testAnalyticsClient) FirstVersionCreated(appSlug, buildSlug, platform string) {
	if c.firstVersionCreatedFn == nil {
		panic("You have to override the FirstVersionCreated function in tests")
	}
	c.firstVersionCreatedFn(appSlug, buildSlug, platform)
}

func (c *testAnalyticsClient) PublishFinished(appSlug string, appVersionID uuid.UUID, result string) {
	if c.publishFinishedFn == nil {
		panic("You have to override the PublishFinished function in tests")
	}
	c.publishFinishedFn(appSlug, appVersionID, result)
}
