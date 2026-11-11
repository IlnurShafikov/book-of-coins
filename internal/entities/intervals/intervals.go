package intervals

// Interval represents a numeric range with an associated identifier.
type Interval struct {
	// ID is the unique identifier for this interval
	ID int

	// Min is the minimum value (inclusive) of the interval
	Min int

	// Max is the maximum value (inclusive) of the interval
	Max int
}

// Intervals is a collection of Interval items.
type Intervals []Interval

// GetIDByValue finds the interval ID that contains the given value.
// Parameters:
//   - value: the value to search for
//
// Returns:
//   - interval ID if found
//   - true if the value falls within an interval, false otherwise
func (v Intervals) GetIDByValue(value int) (int, bool) {
	for _, item := range v {
		if value >= item.Min && value <= item.Max {
			return item.ID, true
		}
	}

	return 0, false
}

// IntervalIDs extracts all interval IDs into a slice.
//
// Returns:
//   - slice of interval IDs in the same order as the intervals
//   - nil if the Intervals collection is nil
func (v Intervals) IntervalIDs() []int {
	if v == nil {
		return nil
	}

	resp := make([]int, len(v))

	for i, interval := range v {
		resp[i] = interval.ID
	}

	return resp
}
