package common

import (
	"fmt"
	"strconv"

	"github.com/IlnurShafikov/book-of-coins/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
	"github.com/IlnurShafikov/slot-engine/pkg/weight"
)

func ReadCoinChanceGroup(
	csvReader *reader.Adapter,
	phasePath *tools.PathBuilder,
) (weight.Group, error) {
	const fileName = "coin_chance_group"

	resp, err := readWeightsGroup(csvReader, phasePath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read final coin values: %w", err)
	}

	return resp, nil
}

// ReadGroup reads and parses generic weight groups from a CSV file.
// This is a general-purpose function for reading weight group configurations.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - phasePath: path builder pointing to the phase configuration directory
//
// Returns:
//   - weights.Group containing weight distributions by group ID
//   - error if reading or parsing fails
func ReadGroup(
	csvReader *reader.Adapter,
	phasePath *tools.PathBuilder,
) (weight.Group, error) {
	const fileName = "weights_group"

	resp, err := readWeightsGroup(csvReader, phasePath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read coin values: %w", err)
	}

	return resp, nil
}

// readWeightsGroup reads and parses weight group data from a CSV file.
// This is a helper function used by other weight group readers.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - filePath: complete path to the CSV file
//
// Returns:
//   - weights.Group containing parsed weight distributions
//   - error if reading or parsing fails
func readWeightsGroup(
	csvReader *reader.Adapter,
	filePath string,
) (weight.Group, error) {
	data, err := readWeightsData(csvReader, filePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", filePath, err)
	}

	group, err := parseWeightsToGroup(data)
	if err != nil {
		return nil, fmt.Errorf("parse weights: %w", err)
	}

	return group, nil
}

// parseWeightsToGroup converts raw CSV data into a weights.Group structure.
// The CSV format expects headers in the first row and values in the first column,
// with weight data in the remaining cells.
//
// Parameters:
//   - data: 2D slice of strings representing CSV data
//
// Returns:
//   - weights.Group containing parsed weight items
//   - error if parsing fails or data format is invalid
func parseWeightsToGroup(data [][]string) (weight.Group, error) {
	const (
		startIndex       = 1
		firstColumnIndex = 0
	)

	err := checkData(data)
	if err != nil {
		return nil, err
	}

	values := data[firstColumnIndex][startIndex:]

	var group weight.Group

	for _, row := range data[startIndex:] {
		groupID, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("parse group id %s: %w", row[0], err)
		}

		items, err := parseWeightItems(values, row[startIndex:])
		if err != nil {
			return nil, fmt.Errorf("parse weight items: %w", err)
		}

		resp, err := weight.NewWeights(items)
		if err != nil {
			return nil, fmt.Errorf("new weights: %w", err)
		}

		group = append(group, weight.NewGroupItem(groupID, resp))
	}

	return group, nil
}

// parseWeightItems creates weight items from value and weight string slices.
// Each item pairs a value with its corresponding weight.
//
// Parameters:
//   - values: slice of value strings
//   - weightsRow: slice of weight strings (must match length of values)
//
// Returns:
//   - slice of weights.Item containing value-weight pairs
//   - error if slices don't match or parsing fails
func parseWeightItems(values, weightsRow []string) ([]weight.Item, error) {
	if len(values) != len(weightsRow) {
		return nil, ErrMismatchedValuesAndWeights
	}

	items := make([]weight.Item, 0, len(values))

	for i, valueStr := range values {
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return nil, fmt.Errorf("parse value: %w", err)
		}

		resp, err := strconv.ParseUint(weightsRow[i], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("parse weight: %w", err)
		}

		items = append(items, weight.NewItem(value, uint32(resp)))
	}

	return items, nil
}

// readWeightsData reads raw CSV data from a file.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - filePath: path to the CSV file
//
// Returns:
//   - 2D slice of strings representing CSV data
//   - error if reading fails
func readWeightsData(
	csvReader *reader.Adapter,
	filePath string,
) ([][]string, error) {
	data, err := csvReader.ReadCSV(filePath)
	if err != nil {
		return nil, fmt.Errorf("read weights data: %w", err)
	}

	return data, nil
}

// checkData validates that the CSV data is not empty.
//
// Parameters:
//   - data: 2D slice of strings representing CSV data
//
// Returns:
//   - ErrEmptyData if data is empty or any row is empty
//   - nil if data is valid
func checkData(data [][]string) error {
	if len(data) == 0 {
		return ErrEmptyData
	}

	for i := range data {
		if len(data[i]) == 0 {
			return ErrEmptyData
		}
	}

	return nil
}
