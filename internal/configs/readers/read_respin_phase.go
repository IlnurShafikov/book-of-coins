package readers

import (
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/components/generator"
	"github.com/IlnurShafikov/book-of-coins/internal/components/manager"
	"github.com/IlnurShafikov/book-of-coins/internal/configs"
	"github.com/IlnurShafikov/book-of-coins/internal/configs/readers/common"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/intervals"
	"github.com/IlnurShafikov/book-of-coins/internal/game/phase/respin"
	"github.com/IlnurShafikov/book-of-coins/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

// readReSpinPhase reads and constructs the respin phase configuration.
// It loads coin value getters, order manager, multiplier getter, and other components
// required for the respin phase.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - modPath: path builder pointing to the game mode directory
//   - general: general configuration shared across all phases
//
// Returns:
//   - pointer to SpinPhase containing respin phase configuration with all components
//   - error if reading or parsing fails
func readReSpinPhase(
	csvReader *reader.Adapter, modPath *tools.PathBuilder, general *configs.GeneralConfig,
) (*respin.SpinPhase, error) {
	const (
		respinPrefix = "respin"
	)

	phasePath := modPath.Phase(respinPrefix)

	gen, err := readRespinGenerator(csvReader, phasePath, general)
	if err != nil {
		return nil, fmt.Errorf("read respin generator: %w", err)
	}

	coinInterval, err := readCoinIntervalsGroup(csvReader, phasePath)
	if err != nil {
		return nil, fmt.Errorf("read coin intervals group: %w", err)
	}

	coinManager := manager.NewCoin(
		general.GrandRules.GrandWin,
		general.GrandRules.GrandJpType,
		general.GrandRules.CoinCountToGrand,
	)

	spinPhase := respin.NewSpinPhase(gen, coinManager, coinInterval)

	return spinPhase, nil
}

// readRespinGenerator builds the respin generator from the coin chance groups and the coin values.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - phasePath: path builder pointing to the respin phase directory
//   - general: general configuration providing the jackpot types
//
// Returns:
//   - pointer to ReSpinGenerator
//   - error if any file cannot be read
func readRespinGenerator(
	csvReader *reader.Adapter,
	phasePath *tools.PathBuilder,
	general *configs.GeneralConfig,
) (*generator.ReSpinGenerator, error) {
	coinChance, err := common.ReadCoinChanceGroup(csvReader, phasePath)
	if err != nil {
		return nil, fmt.Errorf("read coin chance group: %w", err)
	}

	getterValue, err := common.ReadValueGetter(csvReader, phasePath, general)
	if err != nil {
		return nil, fmt.Errorf("read value getter: %w", err)
	}

	gen := generator.NewReSpinGenerator(coinChance, getterValue)

	return gen, nil
}

// readCoinIntervalsGroup reads the intervals that map the number of coins on the board
// to a coin chance group.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - phasePath: path builder pointing to the respin phase directory
//
// Returns:
//   - intervals of coin counts with their coin chance group IDs
//   - error if reading or parsing fails
func readCoinIntervalsGroup(
	csvReader *reader.Adapter,
	phasePath *tools.PathBuilder,
) (intervals.Intervals, error) {
	const fileName = "coin_intervals_group"

	data, err := csvReader.ReadCSV(phasePath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read rs coin intervals: %w", err)
	}

	resp, err := intervals.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse intervals: %w", err)
	}

	return resp, nil
}
