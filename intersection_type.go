package timeline

import "strings"

// IntersectionType defines the possible ways in which two timeline entries can overlap/intersect
type IntersectionType int

const (
	// IntersectionTypeNone indicates that there is no overlap
	IntersectionTypeNone IntersectionType = iota
	// IntersectionTypeSame indicates that the two entries represent the same span of time
	IntersectionTypeSame
	// IntersectionTypeCover indicates that the "new" entry completely covers the "reference" one (i.e. starts earlier and ends later)
	IntersectionTypeCover
	// IntersectionTypeWithin indicates that the "new" entry is completely within the "reference" one (i.e. starts later and ends earlier)
	IntersectionTypeWithin
	// IntersectionTypeAdjacent indicates that the "new" entry is adjacent to the "reference" one
	IntersectionTypeAdjacent
	// IntersectionTypeStartOverlap indicates that the "new" entry overlaps start of the "reference" one (i.e. starts earlier and ends within)
	IntersectionTypeStartOverlap
	// IntersectionTypeEndOverlap indicates that the "new" entry overlaps end of the "reference" one (i.e. starts with and ends later)
	IntersectionTypeEndOverlap
)

// String implements fmt.Stringer for IntersectionType values
func (v IntersectionType) String() string {
	m := map[IntersectionType]string{
		IntersectionTypeNone:         "none",
		IntersectionTypeSame:         "same",
		IntersectionTypeCover:        "cover",
		IntersectionTypeWithin:       "within",
		IntersectionTypeAdjacent:     "adjacent",
		IntersectionTypeStartOverlap: "start",
		IntersectionTypeEndOverlap:   "end",
	}
	if s, ok := m[v]; ok {
		return s
	}
	// this should be unreachable, but IntersectionType(1024) is valid code :(
	return "unknown"
}

// ParseIntersectionType parses the specified string into an IntersectionType enumeration value.
//
// If the string does not contain a valid IntersectionType string, IntersectionTypeNone is returned.
func ParseIntersectionType(s string) IntersectionType {
	v, err := parseIntersectionTypeValue(s)
	if err != nil {
		return IntersectionTypeNone
	}
	return v
}

// MarshalText implements encoding.TextMarshaler for IntersectionType values.
//
// The marshalled value is the result of calling .String() on the enum value.
func (v IntersectionType) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for IntersectionType values.
func (v *IntersectionType) UnmarshalText(p []byte) error {
	r, err := parseIntersectionTypeValue(string(p))
	if err != nil {
		return err
	}
	*v = r
	return nil
}

func parseIntersectionTypeValue(s string) (IntersectionType, error) {
	m := map[string]IntersectionType{
		"none":     IntersectionTypeNone,
		"same":     IntersectionTypeSame,
		"cover":    IntersectionTypeCover,
		"within":   IntersectionTypeWithin,
		"adjacent": IntersectionTypeAdjacent,
		"start":    IntersectionTypeStartOverlap,
		"end":      IntersectionTypeEndOverlap,
	}
	v, exists := m[strings.ToLower(s)]
	if !exists {
		return IntersectionTypeNone, ErrInvalidIntersectionType
	}
	return v, nil
}
