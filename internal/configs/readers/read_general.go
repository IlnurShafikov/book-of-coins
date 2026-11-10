package readers

import (
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/configs"
	"github.com/IlnurShafikov/book-of-coins/internal/configs/readers/common"
	"github.com/IlnurShafikov/book-of-coins/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/paytable"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

const (
	// generalPrefix defines the directory name for general configuration files
	generalPrefix = "general"
	// symbolOnReelName is the CSV file name for the per-reel visible symbol counts.
	symbolOnReelName = "symbols_on_reel"
)

// Read loads and parses the general game configuration.
// It reads the paytable and multipliers from CSV files in the general directory.
//
// Parameters:
//   - reader: CSV file adapter for reading configuration files
//
// Returns:
//   - pointer to GeneralConfig containing paytable and multipliers
//   - error if reading or parsing fails
func Read(reader *reader.Adapter) (*configs.GeneralConfig, error) {
	generalPath := tools.NewPathBuilder(generalPrefix)

	linePaytable, scatterPaytable, err := readPayTables(reader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read paytable: %w", err)
	}

	symbolsOnReel, err := common.ReadIntSlices(reader, generalPath.File(symbolOnReelName))
	if err != nil {
		return nil, fmt.Errorf("read symbols on reel: %w", err)
	}

	winLines, err := common.ReadLines(reader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read win lines: %w", err)
	}

	multipliers, err := common.ReadMultipliers(reader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read multipliers: %w", err)
	}
	fsInitCounter, err := common.ReadFsInitCounter(reader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read fs init counter: %w", err)
	}

	coinIntervals, err := common.ReadCoinIntervals(reader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read coin intervals: %w", err)
	}

	grandRules, err := common.ReadGrandRules(reader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read grand rules: %w", err)
	}

	jpTypes, err := common.ReadJpTypes(reader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read jp types: %w", err)
	}

	expandMinCount, err := common.ReadExpandMinCount(reader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read expand min count: %w", err)
	}

	volatility, err := common.ReadVolatility(reader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read volatility: %w", err)
	}

	return &configs.GeneralConfig{
		PayTable:        linePaytable,
		ScatterPayTable: scatterPaytable,
		Multipliers:     multipliers,
		SymbolsOnReels:  symbolsOnReel,
		WinLines:        winLines,
		JpTypes:         jpTypes,
		ExpandMinCount:  expandMinCount,
		FsInitCounter:   fsInitCounter,
		RsInitCounter:   coinIntervals,
		Volatility:      volatility,
		GrandRules:      grandRules,
	}, nil
}

// readPayTables reads and parses both regular and scatter paytable configurations.
// It loads the main paytable for line wins and a separate paytable for scatter wins.
//
// Parameters:
//   - reader: CSV file adapter for reading configuration files
//   - generalPath: path builder pointing to the general configuration directory
//
// Returns:
//   - paytable: pointer to PayTable for regular line wins
//   - scatterPaytable: pointer to PayTable for scatter wins
//   - err: error if reading or parsing fails
func readPayTables(
	reader *reader.Adapter,
	generalPath *tools.PathBuilder,
) (payTable, scatterPaytable *paytable.PayTable, err error) {
	payTable, err = common.ReadPaytable(reader, generalPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read paytable: %w", err)
	}

	scatterPaytable, err = common.ReadScatterPaytable(reader, generalPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read scatter paytable: %w", err)
	}

	return payTable, scatterPaytable, nil
}
