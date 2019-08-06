package timeline

const (
	// ErrInvalidTimelineStart is returned by NewEntry() if the start time is the zero time
	ErrInvalidTimelineStart = timelineError("The start time must be specified for a timeline entry")
	// ErrInvalidTimelineOrder is returned by NewEntry() if the start time is equal to or later than the end time
	ErrInvalidTimelineOrder = timelineError("The start time must be before the end time for a timeline entry")
	// ErrInvalidIntersectionType indicates that a string could not be parsed into an IntersectionType enum value
	ErrInvalidIntersectionType = timelineError("The provided string could not be parsed into an IntersectionType value")
)

// timelineError defines a custom type so that we can define error constants
type timelineError string

// Error implements error for timelineError values
func (e timelineError) Error() string {
	return string(e)
}
