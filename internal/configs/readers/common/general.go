package common

import (
	"fmt"
	"strconv"

	"github.com/IlnurShafikov/book-of-coins/internal/configs"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/intervals"
	"github.com/IlnurShafikov/book-of-coins/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/paytable"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/fscounter"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

// ReadPaytable reads and parses the paytable configuration from a CSV file.
// The paytable defines symbol payouts and count ranges for winning combinations.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - generalPath: path builder pointing to the general configuration directory
//
// Returns:
//   - pointer to PayTable containing symbol payouts and ranges
//   - error if reading or parsing fails
func ReadPaytable(
	csvReader *reader.Adapter,
	generalPath *tools.PathBuilder,
) (*paytable.PayTable, error) {
	const fileName = "paytable"

	data, err := csvReader.ReadCSV(generalPath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read paytable: %w", err)
	}

	resp, err := paytable.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse paytable: %w", err)
	}

	return resp, nil
}

// ReadScatterPaytable reads and parses the scatter paytable configuration from a CSV file.
// The paytable defines symbol payouts and count ranges for winning combinations.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - generalPath: path builder pointing to the general configuration directory
//
// Returns:
//   - pointer to PayTable containing symbol payouts and ranges
//   - error if reading or parsing fails
func ReadScatterPaytable(
	csvReader *reader.Adapter,
	generalPath *tools.PathBuilder,
) (*paytable.PayTable, error) {
	const fileName = "scatter_paytable"

	data, err := csvReader.ReadCSV(generalPath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read paytable: %w", err)
	}

	resp, err := paytable.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse paytable: %w", err)
	}

	return resp, nil
}

// ReadFsInitCounter reads and parses the free spins initial counter configuration from a CSV file.
// This counter determines the number of free spins awarded at the start.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - generalPath: path builder pointing to the general configuration directory
//
// Returns:
//   - FSInitCounter containing the initial free spins count distribution
//   - error if reading or parsing fails
func ReadFsInitCounter(
	csvReader *reader.Adapter,
	generalPath *tools.PathBuilder,
) (fscounter.FSInitCounter, error) {
	const fileName = "fs_init_counter"

	data, err := csvReader.ReadCSV(generalPath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read data fs init counter: %w", err)
	}

	resp, err := parseDataToSliceInt(data)
	if err != nil {
		return nil, fmt.Errorf("parse data to slice int: %w", err)
	}

	return resp, nil
}

// ReadSymbolsOnReel reads and parses the symbols per reel configuration from a CSV file.
// This configuration defines the number of symbols present on each reel.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - generalPath: path builder pointing to the general configuration directory
//
// Returns:
//   - slice of integers representing the number of symbols on each reel
//   - error if reading or parsing fails
func ReadSymbolsOnReel(
	csvReader *reader.Adapter,
	generalPath *tools.PathBuilder,
) ([]int, error) {
	const fileName = "symbols_on_reel"

	data, err := csvReader.ReadCSV(generalPath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read data symbols on reel: %w", err)
	}

	resp, err := parseDataToSliceInt(data)
	if err != nil {
		return nil, fmt.Errorf("parse data to slice int: %w", err)
	}

	return resp, nil
}

// ReadJpTypes reads and parses the jackpot types configuration from a CSV file.
// This configuration maps jackpot values to their corresponding jackpot type identifiers.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - generalPath: path builder pointing to the general configuration directory
//
// Returns:
//   - map where keys are jackpot values (uint64) and values are jackpot type IDs (int)
//   - error if reading or parsing fails
func ReadJpTypes(
	csvReader *reader.Adapter,
	generalPath *tools.PathBuilder,
) (map[uint64]int, error) {
	const fileName = "jp_types"

	data, err := csvReader.ReadCSV(generalPath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read jp types data: %w", err)
	}

	jpTypes, err := parseJpTypes(data)
	if err != nil {
		return nil, fmt.Errorf("parse jp types: %w", err)
	}

	return jpTypes, nil
}

func ReadExpandMinCount(
	csvReader *reader.Adapter,
	generalPath *tools.PathBuilder,
) (map[int]int, error) {
	const fileName = "expand_min_count"

	data, err := csvReader.ReadCSV(generalPath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read jp types data: %w", err)
	}

	expandMinCount, err := parseExpandMinCount(data)
	if err != nil {
		return nil, fmt.Errorf("parse jp types: %w", err)
	}

	return expandMinCount, nil
}

// ReadCoinIntervals reads and parses the coin count intervals configuration from a CSV file.
// This configuration defines ranges of coin counts mapped to interval IDs for respin logic.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - generalPath: path builder pointing to the general configuration directory
//
// Returns:
//   - Intervals containing coin count ranges with their IDs
//   - error if reading or parsing fails
func ReadCoinIntervals(
	csvReader *reader.Adapter,
	generalPath *tools.PathBuilder,
) (intervals.Intervals, error) {
	const fileName = "rs_coin_intervals"

	data, err := csvReader.ReadCSV(generalPath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read rs coin intervals: %w", err)
	}

	resp, err := intervals.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse intervals: %w", err)
	}

	return resp, nil
}

func ReadGrandRules(
	csvReader *reader.Adapter,
	generalPath *tools.PathBuilder,
) (*configs.GrandRules, error) {
	const (
		coinCountToGrandName = "coin_count_to_grand"
		grandJpTypeName      = "grand_jp_type"
		grandWinName         = "grand_win"
	)

	coinCountToGrand, err := csvReader.ReadInt(generalPath.File(coinCountToGrandName))
	if err != nil {
		return nil, fmt.Errorf("read coin count to grand: %w", err)
	}

	grandJpType, err := csvReader.ReadInt(generalPath.File(grandJpTypeName))
	if err != nil {
		return nil, fmt.Errorf("read grand jackpot type: %w", err)
	}

	grandWin, err := csvReader.ReadUint64(generalPath.File(grandWinName))
	if err != nil {
		return nil, fmt.Errorf("read grand win: %w", err)
	}

	return &configs.GrandRules{
		CoinCountToGrand: coinCountToGrand,
		GrandJpType:      grandJpType,
		GrandWin:         grandWin,
	}, nil
}

// parseJpTypes parses jackpot types configuration from CSV data into a map.
// Each row must contain exactly 2 values: jackpot type ID and jackpot value.
//
// Parameters:
//   - data: 2D array of strings where each row contains [jpType, jpValue]
//
// Returns:
//   - map where keys are jackpot values (uint64) and values are jackpot type IDs (int)
//   - error if data format is invalid or conversion fails
func parseJpTypes(data [][]string) (map[uint64]int, error) {
	const dataLen = 2

	jpTypes := make(map[uint64]int, len(data))

	for _, row := range data {
		if len(row) != dataLen {
			return nil, ErrInvalidData
		}

		jpType, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("parse jp id '%s': %w", row[0], err)
		}

		jpValue, err := strconv.ParseUint(row[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse jp type '%s': %w", row[1], err)
		}

		jpTypes[jpValue] = jpType
	}

	return jpTypes, nil
}

func parseExpandMinCount(data [][]string) (map[int]int, error) {
	const dataLen = 2

	minExpandCount := make(map[int]int, len(data))

	for _, row := range data {
		if len(row) != dataLen {
			return nil, ErrInvalidData
		}

		symbolID, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("parse symbol id '%s': %w", row[0], err)
		}

		minCount, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, fmt.Errorf("parse min count '%s': %w", row[1], err)
		}

		minExpandCount[symbolID] = minCount
	}

	return minExpandCount, nil
}
