package models

import "time"

// TimeService ...
type TimeService struct{}

// Now ...
func (s *TimeService) Now() time.Time {
	return time.Now()
}
