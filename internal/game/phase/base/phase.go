// Package base implements the base game phase logic.
// The base phase represents the regular paid spins in the game.
package base

import (
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/components/calculator"
	"github.com/IlnurShafikov/book-of-coins/internal/components/collectors"
	"github.com/IlnurShafikov/book-of-coins/internal/components/getter"
	"github.com/IlnurShafikov/book-of-coins/internal/components/limiter"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/coins"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/intervals"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/result"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/result/events"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/rules"
	"github.com/IlnurShafikov/book-of-coins/internal/tools/cast"
	"github.com/IlnurShafikov/slot-engine/pkg/generator"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/fscounter"
	"github.com/IlnurShafikov/slot-engine/pkg/weight"
)

// SpinPhase represents the base game spin phase with all its components.
// It contains the game mechanics and configurations needed to execute a base spin.
type SpinPhase struct {
	// gen produces raw reel boards for a given group ID.
	gen *generator.Generator

	getterValues *getter.CoinValue

	// dist define the weighted symbol distributions for reels
	dist *weight.Weights

	// calc performs line, scatter and expand win calculations
	calc calculator.Calculators

	// expandSymbol holds the weighted pool of symbols eligible to expand during free spins
	expandSymbol *weight.Weights

	// fSInitCounter determines initial free spins count
	fSInitCounter fscounter.FSInitCounter

	// rsInitCounter determines initial re spins count
	rsInitCounter intervals.Intervals
}

func NewSpinPhase(
	dist *weight.Weights,
	gen *generator.Generator,
	getterValues *getter.CoinValue,
	calc calculator.Calculators,
	expandSymbol *weight.Weights,
	fsInitCounter fscounter.FSInitCounter,
	rsInitCounter intervals.Intervals,
) *SpinPhase {
	return &SpinPhase{
		dist:          dist,
		gen:           gen,
		getterValues:  getterValues,
		calc:          calc,
		expandSymbol:  expandSymbol,
		fSInitCounter: fsInitCounter,
		rsInitCounter: rsInitCounter,
	}
}

// Spin executes the base game phase and returns the result.
// This is where all base game mechanics are implemented: reel spins, win calculation,
// feature triggers, etc.
//
// Parameters:
//   - rng: random number generator for deterministic randomness
//
// Returns:
//   - pointer to Base containing the spin result
//   - error if spin execution fails
func (b *SpinPhase) Spin(
	rng model.RNG,
	winLimiter limiter.WinLimiter,
) (*result.Base, error) {
	eventer := collectors.NewEvents()
	res := result.NewBase()

	groupID, err := b.dist.GetRandomValue(rng)
	if err != nil {
		return nil, fmt.Errorf("distributions get random value: %w", err)
	}

	board, err := b.gen.Generate(rng, groupID)
	if err != nil {
		return nil, fmt.Errorf("generate board: %w", err)
	}

	totalWin, err := b.processWin(board, eventer, winLimiter)
	if err != nil {
		return nil, fmt.Errorf("process win: %w", err)
	}

	coinResult, err := b.generateCoinValue(rng, board)
	if err != nil {
		return nil, fmt.Errorf("generate coin values: %w", err)
	}

	boardRes := &result.BoardResult{
		Board: board,
		Coins: coinResult,
	}

	res.SetBoardResult(boardRes)

	err = b.checkBonusTrigger(rng, res)
	if err != nil {
		return nil, fmt.Errorf("check bonus trigger: %w", err)
	}

	res.SetTotalWin(totalWin)
	res.SetMaxWinReached(winLimiter.MaxWinReached())
	res.SetEvents(eventer.Get())

	return res, nil
}

// processWin evaluates line wins and coin wins for the generated board result,
// records corresponding events, and returns the combined payout.
// When the win limiter is hit, it returns the capped total immediately.
//
// Parameters:
//   - resp: board result containing the symbol grid and collector state
//   - eventer: event collector that accumulates game events for client delivery
//   - limiter: win cap enforcing the maximum allowed payout
//
// Returns:
//   - combined total of line wins and coin wins, capped by the limiter if reached
//   - error if line calculation or event creation fails
func (b *SpinPhase) processWin(
	board model.Board,
	eventer *collectors.EventsCollector,
	winLimiter limiter.WinLimiter,
) (uint64, error) {
	lineWins, err := b.calc.CalculateLines(board, winLimiter)
	if err != nil {
		return 0, fmt.Errorf("calculate win lines: %w", err)
	}

	scatterWins, err := b.calc.CalculateScatters(board, winLimiter)
	if err != nil {
		return 0, fmt.Errorf("calculate scatter win: %w", err)
	}

	if len(lineWins) > 0 || len(scatterWins) > 0 {
		winsEvent := events.NewWinsEvent(lineWins, scatterWins)
		eventer.Add(winsEvent)
	}

	return winLimiter.TotalWin(), nil
}

// checkBonusTrigger checks if any bonus features (free spins or respins) were triggered.
// It prioritizes free spins by scatter count; if free spins are not triggered,
// it checks respin trigger conditions by coin count.
//
// Parameters:
//   - rng: random number generator used to pick the expand symbol
//   - res: the base result to update with triggered features
//
// Returns:
//   - error if trigger check or configuration fails
func (b *SpinPhase) checkBonusTrigger(rng model.RNG, res *result.Base) error {
	boardResult := res.BoardResult()
	scatterCount := boardResult.Board.CalculateSymbol(rules.ScatterID)
	fsCount := b.fSInitCounter.GetFsInitCount(scatterCount)

	if fsCount > 0 {
		return b.setFreeSpins(rng, res, fsCount)
	}

	coinsCount := boardResult.Board.CalculateSymbol(rules.CoinID)

	if err := b.setReSpins(res, coinsCount); err != nil {
		return fmt.Errorf("set re spin: %w", err)
	}

	return nil
}

// setFreeSpins configures free spins and picks the symbol that will expand during them.
//
// Parameters:
//   - rng: random number generator used to pick the expand symbol
//   - res: the base result to update with free spins configuration
//   - fsCount: number of free spins to grant
//
// Returns:
//   - error if the cast or the expand symbol selection fails
func (b *SpinPhase) setFreeSpins(rng model.RNG, res *result.Base, fsCount int) error {
	safeValue, err := cast.SafeCast[uint32](fsCount)
	if err != nil {
		return fmt.Errorf("safe cast: %w", err)
	}

	expandSymbol, err := b.expandSymbol.GetRandomValue(rng)
	if err != nil {
		return fmt.Errorf("get random value: %w", err)
	}

	res.SetExpandSymbolID(expandSymbol)
	res.SetFreeSpinsInitialsCount(safeValue)

	return nil
}

// setReSpins configures respins if the coin count meets the trigger threshold.
// Parameters:
//   - res: the base result to update with respin configuration
//   - coinCount: number of coins on the board
//
// Returns:
//   - error if respin setup fails
func (b *SpinPhase) setReSpins(
	res *result.Base,
	coinCount int,
) error {
	rsCount, ok := b.rsInitCounter.GetIDByValue(coinCount)
	if !ok || rsCount == 0 {
		return nil
	}

	safeValue, err := cast.SafeCast[uint32](rsCount)
	if err != nil {
		return fmt.Errorf("safe cast: %w", err)
	}

	res.SetReSpinsInitial(safeValue)

	return nil
}

func (b *SpinPhase) generateCoinValue(
	rng model.RNG,
	board model.Board,
) (coins.CoinResult, error) {
	coinCount := board.CalculateSymbol(rules.CoinID)
	if coinCount == 0 {
		return nil, nil
	}

	coinPos := board.SymbolPositions(rules.CoinID)

	coinValues, err := b.getterValues.GenerateValuesInPositions(rng, coinPos)
	if err != nil {
		return nil, fmt.Errorf("generate coin values: %w", err)
	}

	return coinValues, nil
}
