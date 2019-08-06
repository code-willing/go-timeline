package timeline

import (
	"time"

	"github.com/pkg/errors"
)

// New returns a new Timeline consisting of the specified entries
func New(entries ...Entry) Timeline {
	var tl Timeline
	tl.Add(entries...)
	return tl
}

// Timeline represents a slice of Entry instances, sorted by the entries' start time
type Timeline []Entry

// Add adds one or more new entries to an existing timeline and returns a boolean value indicating
// whether or not the timeline was modified
func (tl *Timeline) Add(entries ...Entry) bool {
	updated := false
	for _, e := range entries {
		if tl.addEntry(e) {
			updated = true
		}
	}
	return updated
}

// Normalize sorts the timeline entries by start date and combines any overlapping ranges
//
// This process can be expensive - on the order of O(n^2) - for highly denormalized timeline entries and
// is intended only for cases where there is an existing []Entry that needs to be normalized.
func (tl *Timeline) Normalize() {
	if len(*tl) > 0 {
		// since all of the logic for generating non-overlapping ranges in start date order is already
		// implemented by addEntry(), we can just create a brand new timeline, add each of our entries
		// to it, then replace this timeline w/ the new one wholesale
		ntl := Timeline{}
		for _, e := range *tl {
			ntl.addEntry(e)
		}
		*tl = []Entry(ntl)
	}
}

func (tl *Timeline) addEntry(entry Entry) bool {
	// no work to do if this is the first entry, add it and return
	if len(*tl) == 0 {
		*tl = append(*tl, entry)
		return true
	}
	// step thru existing timeline
	for i, refEntry := range *tl {
		switch Intersect(refEntry, entry) {
		case IntersectionTypeSame, IntersectionTypeWithin:
			// specified date range is already covered, no-op
			return false

		case IntersectionTypeNone:
			// if no intersection and the new entry's start date is before the reference entry, insert at i
			if entry.StartTime().Before(refEntry.StartTime()) {
				*tl = append(*tl, nil)
				copy((*tl)[i+1:], (*tl)[i:])
				(*tl)[i] = entry
				return true
			}

		case IntersectionTypeAdjacent:
			// new entry is adjacent to existing entry
			// . update entry at i w/ new one covering the combined range
			st := entry.StartTime()
			if est := refEntry.StartTime(); est.Before(st) {
				st = est
			}
			et, _ := entry.EndTime()
			if eet, _ := refEntry.EndTime(); eet.After(et) {
				et = eet
			}
			ne, _ := NewEntry(st, et)
			(*tl)[i] = ne
			return true

		case IntersectionTypeStartOverlap:
			// new entry overlaps start of existing entry
			// . update entry at i w/ new one w/ the new start and the existing end
			et, _ := refEntry.EndTime()
			ne, _ := NewEntry(entry.StartTime(), et)
			(*tl)[i] = ne
			return true

		case IntersectionTypeCover, IntersectionTypeEndOverlap:
			// new entry covers or overlaps end of existing entry
			// . update entry at i w/ new one w/ the existing start and the new end
			// . remove any subsequent entries that are covered by the new range
			st := refEntry.StartTime()
			if nst := entry.StartTime(); nst.Before(st) {
				st = nst
			}
			et, _ := entry.EndTime()
			ne, _ := NewEntry(st, et)
			for j, done := i+1, false; !done && j < len(*tl); j++ {
				ee := (*tl)[j]
				itype := Intersect(ee, ne)
				switch itype {
				case IntersectionTypeNone:
					done = true

				case IntersectionTypeCover:
					// remove
					l := len(*tl)
					copy((*tl)[j:], (*tl)[j+1:])
					(*tl)[l-1] = nil
					*tl = (*tl)[:l-1]
					j--

				case IntersectionTypeStartOverlap, IntersectionTypeAdjacent:
					// save the end of this entry
					et, _ = ee.EndTime()
					// remove
					l := len(*tl)
					copy((*tl)[j:], (*tl)[j+1:])
					(*tl)[l-1] = nil
					*tl = (*tl)[:l-1]
					j--

				case IntersectionTypeSame, IntersectionTypeWithin, IntersectionTypeEndOverlap:
					// TODO: remove this panic after everything works
					panic(errors.Errorf("Unexpected intersection type %s at index %d (%s)", itype, j, (*tl)[j]))
				}
			}
			// shift end date of the new entry if we covered >1 existing entry
			if ed, _ := ne.EndTime(); ed.Before(et) {
				ne, _ = NewEntry(st, et)
			}
			(*tl)[i] = ne
			return true
		}
	}
	// new start is later than any existing start w/ no intersection, add to end
	*tl = append(*tl, entry)
	return true
}

// Contains determines whether or not the specified time falls within one of the timeline entries and,
// if it does, returns the start and end of the entry
func (tl Timeline) Contains(t time.Time) (bool, time.Time, time.Time) {
	if t.IsZero() {
		return false, time.Time{}, time.Time{}
	}
	for _, e := range tl {
		sd := e.StartTime()
		ed, hasEnd := e.EndTime()
		if sd.After(t) {
			continue
		}
		if sd.Equal(t) {
			return true, sd, ed
		}
		if !hasEnd || !t.After(ed) {
			return true, sd, ed
		}
	}
	return false, time.Time{}, time.Time{}
}

// Intersect compares two Entry items and returns an IntersectionType value that indicates how the second
// "new" time span intersects with the first "reference" one
func Intersect(refEntry, newEntry Entry) IntersectionType {
	var (
		refStart          = refEntry.StartTime()
		refEnd, refHasEnd = refEntry.EndTime()
		newStart          = newEntry.StartTime()
		newEnd, newHasEnd = newEntry.EndTime()
	)
	if !refHasEnd {
		refEnd = EndOfTime()
	}
	if !newHasEnd {
		newEnd = EndOfTime()
	}
	// newStart < refStart
	if newStart.Before(refStart) {
		if newEnd.Before(refStart) {
			return IntersectionTypeNone
		}
		if newEnd.Equal(refStart) {
			return IntersectionTypeAdjacent
		}
		if newEnd.After(refEnd) {
			return IntersectionTypeCover
		}
		return IntersectionTypeStartOverlap
	}
	// newStart > refStart
	if newStart.After(refStart) {
		if newStart.After(refEnd) {
			return IntersectionTypeNone
		}
		if newStart.Equal(refEnd) {
			return IntersectionTypeAdjacent
		}
		if !newEnd.After(refEnd) {
			return IntersectionTypeWithin
		}
		return IntersectionTypeEndOverlap
	}
	// newStart == refStart
	if newEnd.Before(refEnd) {
		return IntersectionTypeWithin
	}
	if newEnd.After(refEnd) {
		return IntersectionTypeEndOverlap
	}
	// newStart == refStart && newEnd == refEnd
	return IntersectionTypeSame
}
