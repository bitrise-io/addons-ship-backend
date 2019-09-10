package services_test

type testAnalyticsClient struct {
	firstVersionCreatedFn func(appSlug, buildSlug, platform string)
}

func (c *testAnalyticsClient) FirstVersionCreated(appSlug, buildSlug, platform string) {
	if c.firstVersionCreatedFn == nil {
		panic("You have to override the FirstVersionCreated function in tests")
	}
	c.firstVersionCreatedFn(appSlug, buildSlug, platform)
}
