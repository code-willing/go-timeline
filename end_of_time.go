package timeline

import (
	"sync"
	"time"
)

var (
	endOfTimeMtx sync.RWMutex
	endOfTime    = time.Date(9999, time.December, 31, 23, 59, 59, 999999999, time.UTC)
)

// EndOfTime returns the configured sentinel value that represents the latest possible end date for a
// timeline entry, which defaults to midnight UTC on 9999-12-31.
//
// If necessary, this value can be overridden by calling SetEndOfTime()
func EndOfTime() time.Time {
	endOfTimeMtx.RLock()
	v := endOfTime
	endOfTimeMtx.RUnlock()
	return v
}

// SetEndOfTime assigns a custom sentinel value to represent the latest possible end date for a timeline
// entry.
func SetEndOfTime(t time.Time) {
	endOfTimeMtx.Lock()
	endOfTime = t
	endOfTimeMtx.Unlock()
}
