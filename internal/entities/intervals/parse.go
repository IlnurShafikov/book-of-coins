package intervals

import (
	"fmt"
	"strconv"
)

// Parse converts CSV data into an Intervals collection.
// Each row should contain three columns: ID, Min, Max.
//
// Parameters:
//   - data: 2D slice of strings representing CSV data
//
// Returns:
//   - Intervals collection with parsed interval data
//   - error if data is empty, malformed, or validation fails
func Parse(data [][]string) (Intervals, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("%w: empty input", ErrInvalidData)
	}

	result := make(Intervals, len(data))

	for i, row := range data {
		interval, err := parseRow(i, row)
		if err != nil {
			return nil, err
		}

		result[i] = interval
	}

	return result, nil
}

// parseRow parses a single CSV row into an Interval.
//
// Parameters:
//   - rowIndex: index of the row (for error reporting)
//   - row: slice of strings containing ID, Min, Max values
//
// Returns:
//   - parsed Interval
//   - error if row has insufficient columns or invalid values
func parseRow(rowIndex int, row []string) (Interval, error) {
	const (
		idIndex = iota
		minValIndex
		maxValIndex

		expectedCols = 3
	)

	if len(row) < expectedCols {
		return Interval{}, fmt.Errorf("%w: row %d has insufficient columns",
			ErrInvalidData, rowIndex)
	}

	id, err := strconv.Atoi(row[idIndex])
	if err != nil {
		return Interval{}, fmt.Errorf("invalid ID '%s': %w", row[0], err)
	}

	minVal, err := strconv.Atoi(row[minValIndex])
	if err != nil {
		return Interval{}, fmt.Errorf("invalid Min '%s': %w", row[1], err)
	}

	maxVal, err := strconv.Atoi(row[maxValIndex])
	if err != nil {
		return Interval{}, fmt.Errorf("invalid Max '%s': %w", row[2], err)
	}

	if err := validateIntervalValue(minVal, maxVal); err != nil {
		return Interval{}, fmt.Errorf("row %d: %w", rowIndex, err)
	}

	return Interval{
		ID:  id,
		Min: minVal,
		Max: maxVal,
	}, nil
}

// validateIntervalValue checks if interval min and max values are valid.
// Validation rules:
//   - minVal must be non-negative
//   - maxVal must be positive
//   - minVal must be less than maxVal
//   - minVal and maxVal must not be equal
//
// Parameters:
//   - minVal: minimum value of the interval
//   - maxVal: maximum value of the interval
//
// Returns:
//   - error if validation fails, nil otherwise
func validateIntervalValue(minVal, maxVal int) error {
	if minVal < 0 {
		return fmt.Errorf("%w: %d", ErrInvalidMin, minVal)
	}

	if maxVal <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidMax, maxVal)
	}

	if minVal > maxVal {
		return fmt.Errorf("%w: min=%d max=%d", ErrMinGreaterThanMax, minVal, maxVal)
	}

	if minVal == maxVal {
		return fmt.Errorf("%w: %d", ErrMinEqualMax, minVal)
	}

	return nil
}
