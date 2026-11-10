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
//   - reader: CSV file adapter for reading configuration files
//   - modPath: path builder pointing to the game mode directory
//   - general: general configuration shared across all phases
//
// Returns:
//   - pointer to SpinPhase containing respin phase configuration with all components
//   - error if reading or parsing fails
func readReSpinPhase(
	reader *reader.Adapter, modPath *tools.PathBuilder, general *configs.GeneralConfig,
) (*respin.SpinPhase, error) {
	const (
		respinPrefix = "respin"
	)

	phasePath := modPath.Phase(respinPrefix)

	gen, err := readRespinGenerator(reader, phasePath, general)
	if err != nil {
		return nil, fmt.Errorf("read respin generator: %w", err)
	}

	coinInterval, err := readCoinIntervalsGroup(reader, phasePath)
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

func readRespinGenerator(
	reader *reader.Adapter,
	phasePath *tools.PathBuilder,
	general *configs.GeneralConfig,
) (*generator.ReSpinGenerator, error) {
	coinChance, err := common.ReadCoinChanceGroup(reader, phasePath)
	if err != nil {
		return nil, fmt.Errorf("read coin chance group: %w", err)
	}

	getterValue, err := common.ReadValueGetter(reader, phasePath, general)
	if err != nil {
		return nil, fmt.Errorf("read value getter: %w", err)
	}

	gen := generator.NewReSpinGenerator(coinChance, getterValue)

	return gen, nil
}
func readCoinIntervalsGroup(
	reader *reader.Adapter,
	phasePath *tools.PathBuilder,
) (intervals.Intervals, error) {
	const fileName = "coin_intervals_group"

	data, err := reader.ReadCSV(phasePath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read rs coin intervals: %w", err)
	}

	resp, err := intervals.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse intervals: %w", err)
	}

	return resp, nil
}
