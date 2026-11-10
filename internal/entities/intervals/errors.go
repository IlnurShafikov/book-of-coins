// Package intervals provides data structures and functions for working with numeric intervals.
// Intervals are used to categorize numeric values into ranges with associated IDs.
package intervals

import "errors"

var (
	// ErrInvalidData is returned when the input data is empty or has invalid format.
	ErrInvalidData = errors.New("invalid data")

	// ErrInvalidMin is returned when the minimum value is negative.
	ErrInvalidMin = errors.New("min value is invalid")

	// ErrInvalidMax is returned when the maximum value is zero or negative.
	ErrInvalidMax = errors.New("max value is invalid")

	// ErrMinGreaterThanMax is returned when the minimum value exceeds the maximum value.
	ErrMinGreaterThanMax = errors.New("min is greater than max")

	// ErrMinEqualMax is returned when the minimum and maximum values are equal.
	ErrMinEqualMax = errors.New("min equals max")
)
