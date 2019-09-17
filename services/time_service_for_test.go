package services_test

import "time"

type testTimeService struct {
	nowFn func() time.Time
}

func (c *testTimeService) Now() time.Time {
	if c.nowFn == nil {
		panic("You have to override the Now function in tests")
	}

	return c.nowFn()
}
