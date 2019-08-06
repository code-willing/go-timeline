package timeline

import (
	"fmt"
	"time"
)

// Entry defines a type that represents a single entry in a timeline with a start time and, optionally,
// an end time that is after the start time.
type Entry interface {
	StartTime() time.Time
	EndTime() (time.Time, bool)
	Duration() time.Duration
}

// Must panics if err is non-nil, otherwise it returns e
func Must(e Entry, err error) Entry {
	if err != nil {
		panic(err)
	}
	return e
}

// NewEntry creates a new timeline entry with the specified start and end dates.
//
// If end is the zero time, this entry will have no end date
func NewEntry(start time.Time, end time.Time) (Entry, error) {
	if start.IsZero() {
		return nil, ErrInvalidTimelineStart
	}
	if !(end.IsZero() || start.Before(end)) {
		return nil, ErrInvalidTimelineOrder
	}
	if end.IsZero() {
		end = EndOfTime()
	}
	e := entry{
		start: start,
		end:   end,
	}
	return e, nil
}

// FromStartDate creates a new timeline entry starting on the specified date (at midnight UTC) and no end date
func FromStartDate(y int, m time.Month, d int) (Entry, error) {
	return NewEntry(
		time.Date(y, m, d, 0, 0, 0, 0, time.UTC),
		time.Time{},
	)
}

// ForDateRange creates a new timeline entry starting and ending at midnight UTC on the specified dates
func ForDateRange(sy int, sm time.Month, sd int, ey int, em time.Month, ed int) (Entry, error) {
	return NewEntry(
		time.Date(sy, sm, sd, 0, 0, 0, 0, time.UTC),
		time.Date(ey, em, ed, 0, 0, 0, 0, time.UTC),
	)
}

type entry struct {
	start time.Time
	end   time.Time
}

// StartTime implements Entry.StartTime() and returns the start time for this timeline entry
func (e entry) StartTime() time.Time {
	return e.start
}

// EndTime implements Entry.EndTime() and returns the end time for this timeline entry, along with a
// boolean value indicating if an end time was present
func (e entry) EndTime() (time.Time, bool) {
	return e.end, !e.end.Equal(EndOfTime())
}

// Duration returns the period between the start and end time for this timeline entry.  If the entry
// does not have an end, the period between the start and EndOfTime (midnight UTC on 9999-12-31) is returned
func (e entry) Duration() time.Duration {
	return e.end.Sub(e.start)
}

// String implements fmt.Stringer for entry instances
//
// The returned string contains the span of the entry in range notation with the following format: [<start> .. <end>].
// If this entry does not have an end date, the result is formatted as "[<start> .. -)" (to indicate no upper bound).
// The date values are printed according to the time.RFC3339 format.
func (e entry) String() string {
	if end, hasEnd := e.EndTime(); hasEnd {
		return fmt.Sprintf("[%s .. %s]", e.start.Format(time.RFC3339), end.Format(time.RFC3339))
	}
	return fmt.Sprintf("[%s .. -)", e.start.Format(time.RFC3339))
}
