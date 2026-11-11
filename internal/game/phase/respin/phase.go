// Package respin implements the respin phase of the slot game.
// During respins only coins stay on the board; every respin may drop new coins
// into empty cells, and each new coin restores the respins counter.
package respin

import (
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/components/collectors"
	"github.com/IlnurShafikov/book-of-coins/internal/components/generator"
	"github.com/IlnurShafikov/book-of-coins/internal/components/limiter"
	"github.com/IlnurShafikov/book-of-coins/internal/components/manager"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/intervals"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/result"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/rules"
	"github.com/IlnurShafikov/book-of-coins/internal/tools/cast"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
)

// SpinPhase represents the respin phase with all its components.
type SpinPhase struct {
	// gen places new coins on the empty cells of the board
	gen *generator.ReSpinGenerator

	// coinManager collects the coins and pays out the accumulated win
	coinManager *manager.Coin

	// coinIntervalGroupID maps the current coin count to a coin chance group
	coinIntervalGroupID intervals.Intervals
}

// NewSpinPhase creates a new respin phase with all required components.
//
// Parameters:
//   - gen: generator that drops new coins during respins
//   - coinManager: collector that calculates the final coin win
//   - coinIntervalGroupID: intervals mapping the coin count to a chance group
//
// Returns:
//   - pointer to a new SpinPhase instance
func NewSpinPhase(
	gen *generator.ReSpinGenerator,
	coinManager *manager.Coin,
	coinIntervalGroupID intervals.Intervals,
) *SpinPhase {
	return &SpinPhase{
		gen:                 gen,
		coinManager:         coinManager,
		coinIntervalGroupID: coinIntervalGroupID,
	}
}

// Spins plays the whole respin session until the respins counter is exhausted
// or the board is full. The win is paid once, on the last respin.
//
// The method modifies prevBoardResult in place (non-coin symbols are cleared and
// new coins are added), so the caller must pass a copy of the board it wants to keep.
//
// Parameters:
//   - rng: random number generator
//   - rsCount: initial number of respins
//   - prevBoardResult: board of the phase that triggered the respins (modified)
//   - winLimiter: limiter enforcing maximum win cap
//
// Returns:
//   - results of all respins in the order they were played
//   - error if generation or counter calculation fails
func (r *SpinPhase) Spins(
	rng model.RNG,
	rsCount int,
	prevBoardResult *result.BoardResult,
	winLimiter limiter.WinLimiter,
) ([]*result.ReSpin, error) {
	const avrLengthSession = 10

	clearBoard(prevBoardResult.Board)

	allRes := make([]*result.ReSpin, 0, avrLengthSession)
	coinCount := prevBoardResult.Board.CalculateSymbol(rules.CoinID)

	for rsCount > 0 {
		res, left, err := r.spin(rng, rsCount, coinCount, prevBoardResult, winLimiter)
		if err != nil {
			return nil, err
		}

		coinCount = res.BoardResult().Board.CalculateSymbol(rules.CoinID)
		rsCount = left

		allRes = append(allRes, res)
	}

	return allRes, nil
}

// spin plays a single respin: drops new coins, applies the retrigger and, if it was
// the last respin, collects the coins.
//
// Parameters:
//   - rng: random number generator
//   - rsCount: respins left before this spin
//   - coinCount: coins on the board before this spin
//   - board: current board (modified: new coins are placed on it)
//   - winLimiter: limiter enforcing maximum win cap
//
// Returns:
//   - result of this respin (its board is the state after new coins were placed)
//   - respins left after this spin
//   - error if generation or counter calculation fails
func (r *SpinPhase) spin(
	rng model.RNG,
	rsCount, coinCount int,
	board *result.BoardResult,
	winLimiter limiter.WinLimiter,
) (res *result.ReSpin, left int, err error) {
	eventer := collectors.NewEvents()
	res = result.NewReSpin()
	left = rsCount - 1

	groupID, ok := r.coinIntervalGroupID.GetIDByValue(coinCount)
	if !ok {
		return nil, 0, ErrInvalidIntervalCoinCount
	}

	addedCoin, err := r.gen.Generate(rng, groupID, board)
	if err != nil {
		return nil, 0, fmt.Errorf("generate new coin: %w", err)
	}

	res.SetBoardResult(board.Copy())

	if addedCoin > 0 {
		retrigger, rtErr := handleReTrigger(left)
		if rtErr != nil {
			return nil, 0, fmt.Errorf("handle retrigger: %w", rtErr)
		}

		res.SetReTriggers(retrigger)
	}

	left = r.applyRetrigger(left, res)

	if left == 0 {
		res.SetStepWin(r.coinManager.Collect(res.BoardResult(), eventer, winLimiter))
	}

	res.SetEvents(eventer.Get())
	res.SetTotalWin(winLimiter.TotalWin())
	res.SetMaxWinReached(winLimiter.MaxWinReached())

	return res, left, nil
}

// applyRetrigger finishes the session when the board is full and restores the
// respins counter when the current respin dropped a new coin.
//
// Parameters:
//   - rsCount: respins left after the current spin
//   - res: result of the current respin (retriggers are reset if the board is full)
//
// Returns:
//   - respins left after the retrigger rules are applied
func (r *SpinPhase) applyRetrigger(rsCount int, res *result.ReSpin) int {
	voidCount := res.BoardResult().Board.CalculateSymbol(model.Void)

	if voidCount == 0 {
		res.SetReTriggers(0)

		return 0
	}

	if res.ReTriggers() > 0 {
		return rules.InitRespinsCount
	}

	return rsCount
}

// handleReTrigger calculates how many respins are added to restore the counter.
//
// Parameters:
//   - rsCount: respins left after the current spin
//
// Returns:
//   - number of respins added
//   - error if the value cannot be safely cast
func handleReTrigger(rsCount int) (uint32, error) {
	safeValue, err := cast.SafeCast[uint32](rules.InitRespinsCount - rsCount)
	if err != nil {
		return 0, fmt.Errorf("safe cast value: %w", err)
	}

	return safeValue, nil
}

// clearBoard removes all non-coin symbols from the board, replacing them with Void.
// This prepares the board for the respin phase where only coins persist.
//
// Parameters:
//   - board: the game board to clear
func clearBoard(board model.Board) {
	for reel, symbols := range board {
		for row, symbolID := range symbols {
			if symbolID != rules.CoinID {
				board[reel][row] = model.Void
			}
		}
	}
}
