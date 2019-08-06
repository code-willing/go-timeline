package timeline_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/code-willing/timeline"
)

func TestIntersect(t *testing.T) {
	cases := []struct {
		name     string
		v1       timeline.Entry
		v2       timeline.Entry
		expected timeline.IntersectionType
	}{
		{
			"same values",
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.IntersectionTypeSame,
		},
		{
			"new completely before ref",
			timeline.Must(timeline.NewEntry(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.IntersectionTypeNone,
		},
		{
			"new end same as ref start",
			timeline.Must(timeline.NewEntry(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))),
			timeline.IntersectionTypeAdjacent,
		},
		{
			"earlier start/end between ref start and end",
			timeline.Must(timeline.NewEntry(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.February, 1, 0, 0, 0, 0, time.UTC))),
			timeline.IntersectionTypeStartOverlap,
		},
		{
			"earlier start/end same as ref end",
			timeline.Must(timeline.NewEntry(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.IntersectionTypeStartOverlap,
		},
		{
			"earlier start/later end",
			timeline.Must(timeline.NewEntry(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC))),
			timeline.IntersectionTypeCover,
		},
		{
			"same start/earlier end",
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(1999, time.December, 30, 0, 0, 0, 0, time.UTC))),
			timeline.IntersectionTypeWithin,
		},
		{
			"same start/later end",
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))),
			timeline.IntersectionTypeEndOverlap,
		},
		{
			"start within ref range/end before ref end",
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2001, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.Must(timeline.NewEntry(time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC), time.Date(2000, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.IntersectionTypeWithin,
		},
		{
			"start within ref range/end same as ref end",
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2001, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.Must(timeline.NewEntry(time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC), time.Date(2001, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.IntersectionTypeWithin,
		},
		{
			"start within ref range/end after ref end",
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2001, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.Must(timeline.NewEntry(time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC), time.Date(2002, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.IntersectionTypeEndOverlap,
		},
		{
			"start same as ref end/later end",
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.Must(timeline.NewEntry(time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC), time.Date(2000, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.IntersectionTypeAdjacent,
		},
		{
			"new completely after ref value",
			timeline.Must(timeline.NewEntry(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.Must(timeline.NewEntry(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.December, 31, 0, 0, 0, 0, time.UTC))),
			timeline.IntersectionTypeNone,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			got := timeline.Intersect(tc.v1, tc.v2)
			if got != tc.expected {
				tt.Errorf("Expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestAddEntry(t *testing.T) {
	cases := []struct {
		name         string
		value        timeline.Timeline
		newEntries   []timeline.Entry
		expected     timeline.Timeline
		shouldChange bool
	}{
		{
			"empty timeline/single new entry",
			timeline.New(),
			[]timeline.Entry{
				timeline.Must(timeline.FromStartDate(2000, time.January, 1)),
			},
			timeline.New(
				timeline.Must(timeline.FromStartDate(2000, time.January, 1)),
			),
			true,
		},
		{
			"empty timeline/non-contiguous new entries",
			timeline.New(),
			[]timeline.Entry{
				timeline.Must(timeline.ForDateRange(1998, time.January, 1, 1999, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2000, time.January, 1)),
			},
			timeline.New(
				timeline.Must(timeline.ForDateRange(1998, time.January, 1, 1999, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2000, time.January, 1)),
			),
			true,
		},
		{
			"empty timeline/contiguous new entries",
			timeline.New(),
			[]timeline.Entry{
				timeline.Must(timeline.ForDateRange(1998, time.January, 1, 2000, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2000, time.January, 1)),
			},
			timeline.New(
				timeline.Must(timeline.FromStartDate(1998, time.January, 1)),
			),
			true,
		},
		{
			"empty timeline/overlapping new entries/new after existing",
			timeline.New(),
			[]timeline.Entry{
				timeline.Must(timeline.ForDateRange(1998, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2000, time.January, 1)),
			},
			timeline.New(
				timeline.Must(timeline.FromStartDate(1998, time.January, 1)),
			),
			true,
		},
		{
			"empty timeline/overlapping new entries/new before existing",
			timeline.New(),
			[]timeline.Entry{
				timeline.Must(timeline.FromStartDate(2000, time.January, 1)),
				timeline.Must(timeline.ForDateRange(1998, time.January, 1, 2001, time.January, 1)),
			},
			timeline.New(
				timeline.Must(timeline.FromStartDate(1998, time.January, 1)),
			),
			true,
		},
		{
			"empty timeline/overlapping new entries/new covers existing",
			timeline.New(),
			[]timeline.Entry{
				timeline.Must(timeline.ForDateRange(1998, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.FromStartDate(1997, time.January, 1)),
			},
			timeline.New(
				timeline.Must(timeline.FromStartDate(1997, time.January, 1)),
			),
			true,
		},
		{
			"empty timeline/overlapping new entries/new within existing",
			timeline.New(),
			[]timeline.Entry{
				timeline.Must(timeline.FromStartDate(1997, time.January, 1)),
				timeline.Must(timeline.ForDateRange(1998, time.January, 1, 2001, time.January, 1)),
			},
			timeline.New(
				timeline.Must(timeline.FromStartDate(1997, time.January, 1)),
			),
			true,
		},
		{
			"existing timeline/non-overlapping new entry",
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			[]timeline.Entry{
				timeline.Must(timeline.ForDateRange(2002, time.January, 1, 2003, time.January, 1)),
			},
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2002, time.January, 1, 2003, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			true,
		},
		{
			"existing timeline/new entry matches existing",
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			[]timeline.Entry{
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
			},
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			false,
		},
		{
			"existing timeline/contiguous new entry/extend end of existing",
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			[]timeline.Entry{
				timeline.Must(timeline.ForDateRange(2001, time.January, 1, 2002, time.January, 1)),
			},
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2002, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			true,
		},
		{
			"existing timeline/contiguous new entry/extend start of existing backwards",
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			[]timeline.Entry{
				timeline.Must(timeline.ForDateRange(2002, time.January, 1, 2004, time.January, 1)),
			},
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2002, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			true,
		},
		{
			"existing timeline/overlapping new entry/extend start of existing backwards",
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			[]timeline.Entry{
				timeline.Must(timeline.ForDateRange(2002, time.January, 1, 2004, time.June, 1)),
			},
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2002, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			true,
		},
		{
			"existing timeline/overlapping new entry/extend end of existing",
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			[]timeline.Entry{
				timeline.Must(timeline.ForDateRange(2004, time.June, 1, 2006, time.January, 1)),
			},
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2006, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			true,
		},
		{
			"existing timeline/overlapping new entry/cover existing",
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			[]timeline.Entry{
				timeline.Must(timeline.ForDateRange(2003, time.January, 1, 2006, time.January, 1)),
			},
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2003, time.January, 1, 2006, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			true,
		},
		{
			"existing timeline/overlapping new entry/within existing",
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			[]timeline.Entry{
				timeline.Must(timeline.ForDateRange(2004, time.January, 2, 2004, time.December, 31)),
			},
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			false,
		},
		{
			"existing timeline/overlapping new entry/merge multiple existing",
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			[]timeline.Entry{
				timeline.Must(timeline.ForDateRange(2000, time.June, 1, 2004, time.June, 1)),
			},
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			true,
		},
		{
			"existing timeline/overlapping new entry/merge all existing",
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.FromStartDate(2010, time.January, 1)),
			),
			[]timeline.Entry{
				timeline.Must(timeline.ForDateRange(2000, time.June, 1, 2010, time.June, 1)),
			},
			timeline.New(
				timeline.Must(timeline.FromStartDate(2000, time.January, 1)),
			),
			true,
		},
		{
			"existing timeline/overlapping new entry/end past all existing",
			timeline.New(
				timeline.Must(timeline.ForDateRange(2000, time.January, 1, 2001, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2004, time.January, 1, 2005, time.January, 1)),
				timeline.Must(timeline.ForDateRange(2010, time.January, 1, 2011, time.January, 1)),
			),
			[]timeline.Entry{
				timeline.Must(timeline.FromStartDate(2000, time.June, 1)),
			},
			timeline.New(
				timeline.Must(timeline.FromStartDate(2000, time.January, 1)),
			),
			true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			changed := tc.value.Add(tc.newEntries...)
			if changed != tc.shouldChange {
				tt.Errorf("Expected 'changed' to be:\n\t%v\nGot:\n\t%v", tc.shouldChange, changed)
			}
			if !testIsSameTimeline(tc.value, tc.expected) {
				tt.Errorf("Expected:\n\t%s\nGot:\n\t%s", printTimeline(tc.expected), printTimeline(tc.value))
			}
		})
	}
}

func testIsSameTimeline(tl1, tl2 timeline.Timeline) bool {
	if len(tl1) != len(tl2) {
		return false
	}
	for i, e := range tl1 {
		e2 := tl2[i]
		if !testIsSameEntry(e, e2) {
			return false
		}
	}
	return true
}

func testIsSameEntry(e1, e2 timeline.Entry) bool {
	if !e1.StartTime().Equal(e2.StartTime()) {
		return false
	}
	end1, _ := e1.EndTime()
	end2, _ := e2.EndTime()
	return end1.Equal(end2)
}

func printTimeline(tl timeline.Timeline) string {
	return fmt.Sprintf("%v", tl)
}
