package dataservices

import "time"

// TimeInterface ...
type TimeInterface interface {
	Now() time.Time
}
