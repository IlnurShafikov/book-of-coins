package readers

import (
	"fmt"

	calc "github.com/IlnurShafikov/book-of-coins/internal/components/calculator"
	"github.com/IlnurShafikov/book-of-coins/internal/configs"
	"github.com/IlnurShafikov/book-of-coins/internal/configs/readers/common"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/rules"
	"github.com/IlnurShafikov/book-of-coins/internal/game/phase/base"
	"github.com/IlnurShafikov/book-of-coins/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/calculators"
	"github.com/IlnurShafikov/slot-engine/pkg/generator"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

// readBasePhase reads and constructs the base game phase configuration.
// It loads symbol distributions and weight groups from CSV files.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - modPath: path builder pointing to the game mode directory
//   - general: general configuration shared across all phases
//
// Returns:
//   - pointer to SpinPhase containing base game configuration
//   - error if reading or parsing fails
func readBasePhase(
	csvReader *reader.Adapter, modPath *tools.PathBuilder, general *configs.GeneralConfig,
) (*base.SpinPhase, error) {
	const (
		basePrefix       = "base"
		expandSymbolName = "expand_symbol"
	)

	// PathBuilder for creating paths to files
	phasePath := modPath.Phase(basePrefix)

	dist, err := common.ReadDistributions(csvReader, phasePath)
	if err != nil {
		return nil, fmt.Errorf("read distributions: %w", err)
	}

	gen, err := readReelsGenerator(csvReader, phasePath, general, dist.CopyValues())
	if err != nil {
		return nil, fmt.Errorf("read base reel generator: %w", err)
	}

	phaseCalc := readCalculators(general)

	expandSymbol, err := common.ReadWeight(csvReader, phasePath.File(expandSymbolName))
	if err != nil {
		return nil, fmt.Errorf("read expand symbol weight: %w", err)
	}

	getterValue, err := common.ReadValueGetter(csvReader, phasePath, general)
	if err != nil {
		return nil, fmt.Errorf("read value getter: %w", err)
	}

	spinPhase := base.NewSpinPhase(
		dist, gen, getterValue, phaseCalc, expandSymbol,
		general.FsInitCounter, general.RsInitCounter,
	)

	return spinPhase, nil
}

// readCalculators creates a line-based win calculator for the base game.
// Configures wild symbols and symbols to skip during line win calculations.
//
// Parameters:
//   - general: general configuration containing paytable and win lines
//
// Returns:
//   - pointer to Calculator configured for line-based wins
func readCalculators(
	general *configs.GeneralConfig,
) *calc.Calculator {
	wild := model.Collection{rules.ScatterID}
	scatter := model.Collection{rules.ScatterID}
	skipSymbol := model.Collection{rules.CoinID}

	linesCalc := calculators.NewLines(
		general.PayTable,
		general.WinLines,
		wild,
		skipSymbol)

	scatterCalc := calculators.NewPayAnyWhere(general.ScatterPayTable, scatter)
	expandCalc := calc.NewExpandCalculator(general.PayTable, general.WinLines)

	return calc.NewCalculators(linesCalc, scatterCalc, expandCalc)
}

// readReelsGenerator reads reel-strip data for every trigger group and
// builds the standard board generator used for base-game spins.
//
// Parameters:
//   - csvReader: CSV file adapter backed by the embedded config FS
//   - phasePath: path builder pointing to the base-game phase directory
//   - general: shared game configuration supplying the per-reel symbol counts
//   - groupIDs: trigger group IDs for which reel files must be loaded
//
// Returns:
//   - pointer to Generator with all group reel strips loaded
//   - error if any reel file cannot be read or the generator cannot be initialised
func readReelsGenerator(
	csvReader *reader.Adapter,
	phasePath *tools.PathBuilder,
	general *configs.GeneralConfig,
	groupIDs []model.GroupID,
) (*generator.Generator, error) {
	reels, err := common.ReadReels(csvReader, phasePath, groupIDs)
	if err != nil {
		return nil, fmt.Errorf("read reels: %w", err)
	}

	reelGen, err := generator.NewReelsGenerator(reels, general.SymbolsOnReels)
	if err != nil {
		return nil, fmt.Errorf("make base board generator: %w", err)
	}

	return reelGen, nil
}
