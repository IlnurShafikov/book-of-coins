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
//   - csvReader: CSV file adapter for reading configuration files
//
// Returns:
//   - pointer to GeneralConfig containing paytable and multipliers
//   - error if reading or parsing fails
func Read(csvReader *reader.Adapter) (*configs.GeneralConfig, error) {
	generalPath := tools.NewPathBuilder(generalPrefix)

	linePaytable, scatterPaytable, err := readPayTables(csvReader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read paytable: %w", err)
	}

	symbolsOnReel, err := common.ReadIntSlices(csvReader, generalPath.File(symbolOnReelName))
	if err != nil {
		return nil, fmt.Errorf("read symbols on reel: %w", err)
	}

	winLines, err := common.ReadLines(csvReader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read win lines: %w", err)
	}

	multipliers, err := common.ReadMultipliers(csvReader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read multipliers: %w", err)
	}

	fsInitCounter, err := common.ReadFsInitCounter(csvReader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read fs init counter: %w", err)
	}

	expandMinCount, err := common.ReadExpandMinCount(csvReader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read expand min count: %w", err)
	}

	volatility, err := common.ReadVolatility(csvReader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read volatility: %w", err)
	}

	cfg := &configs.GeneralConfig{
		PayTable:        linePaytable,
		ScatterPayTable: scatterPaytable,
		Multipliers:     multipliers,
		SymbolsOnReels:  symbolsOnReel,
		WinLines:        winLines,
		JpTypes:         nil,
		ExpandMinCount:  expandMinCount,
		FsInitCounter:   fsInitCounter,
		RsInitCounter:   nil,
		Volatility:      volatility,
		GrandRules:      nil,
	}

	if err = readCoinConfigs(csvReader, generalPath, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// readCoinConfigs reads coin related settings (respin intervals, grand rules and
// jackpot types) and stores them in cfg.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - generalPath: path builder pointing to the general configuration directory
//   - cfg: general configuration to fill
//
// Returns:
//   - error if reading or parsing fails
func readCoinConfigs(
	csvReader *reader.Adapter,
	generalPath *tools.PathBuilder,
	cfg *configs.GeneralConfig,
) (err error) {
	cfg.RsInitCounter, err = common.ReadCoinIntervals(csvReader, generalPath)
	if err != nil {
		return fmt.Errorf("read coin intervals: %w", err)
	}

	cfg.GrandRules, err = common.ReadGrandRules(csvReader, generalPath)
	if err != nil {
		return fmt.Errorf("read grand rules: %w", err)
	}

	cfg.JpTypes, err = common.ReadJpTypes(csvReader, generalPath)
	if err != nil {
		return fmt.Errorf("read jp types: %w", err)
	}

	return nil
}

// readPayTables reads and parses both regular and scatter paytable configurations.
// It loads the main paytable for line wins and a separate paytable for scatter wins.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - generalPath: path builder pointing to the general configuration directory
//
// Returns:
//   - paytable: pointer to PayTable for regular line wins
//   - scatterPaytable: pointer to PayTable for scatter wins
//   - err: error if reading or parsing fails
func readPayTables(
	csvReader *reader.Adapter,
	generalPath *tools.PathBuilder,
) (payTable, scatterPaytable *paytable.PayTable, err error) {
	payTable, err = common.ReadPaytable(csvReader, generalPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read paytable: %w", err)
	}

	scatterPaytable, err = common.ReadScatterPaytable(csvReader, generalPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read scatter paytable: %w", err)
	}

	return payTable, scatterPaytable, nil
}
