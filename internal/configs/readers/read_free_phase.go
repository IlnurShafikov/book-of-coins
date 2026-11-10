package readers

import (
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/configs"
	"github.com/IlnurShafikov/book-of-coins/internal/configs/readers/common"
	"github.com/IlnurShafikov/book-of-coins/internal/game/phase/free"
	"github.com/IlnurShafikov/book-of-coins/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

// readFreePhase reads and constructs the free spins phase configuration.
// Currently returns a basic configuration; extend as needed for your game.
//
// Parameters:
//   - reader: CSV file adapter for reading configuration files
//   - modPath: path builder pointing to the game mode directory
//   - general: general configuration shared across all phases
//
// Returns:
//   - pointer to SpinPhase containing free spins configuration
//   - error if reading or parsing fails
func readFreePhase(
	reader *reader.Adapter,
	modPath *tools.PathBuilder,
	general *configs.GeneralConfig,
) (*free.SpinPhase, error) {
	const (
		freePrefix = "free"
	)

	phasePath := modPath.Phase(freePrefix)

	dist, err := common.ReadDistributions(reader, phasePath)
	if err != nil {
		return nil, fmt.Errorf("read distributions: %w", err)
	}

	reelsGenerator, err := readReelsGenerator(reader, phasePath, general, dist.CopyValues())
	if err != nil {
		return nil, fmt.Errorf("read reels generator: %w", err)
	}

	calc := readCalculators(general)

	getterValue, err := common.ReadValueGetter(reader, phasePath, general)
	if err != nil {
		return nil, fmt.Errorf("read value getter: %w", err)
	}

	spinPhase := free.NewSpinPhase(
		reelsGenerator,
		calc,
		getterValue,
		general.ExpandMinCount,
		general.FsInitCounter,
		general.RsInitCounter,
	)

	return spinPhase, nil
}
