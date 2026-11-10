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

type SpinPhase struct {
	gen                 *generator.ReSpinGenerator
	coinManager         *manager.Coin
	coinIntervalGroupID intervals.Intervals
}

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

func (r *SpinPhase) Spins(
	rng model.RNG,
	rsCount int,
	prevBoardResult *result.BoardResult,
	winLimiter limiter.WinLimiter,
) ([]*result.ReSpin, error) {
	const avrLenghtSession = 10

	clearBoard(prevBoardResult.Board)

	allRes := make([]*result.ReSpin, 0, avrLenghtSession)
	coinCount := prevBoardResult.Board.CalculateSymbol(rules.CoinID)

	for rsCount > 0 {
		eventer := collectors.NewEvents()
		res := result.NewReSpin()

		rsCount--

		groupID, ok := r.coinIntervalGroupID.GetIDByValue(coinCount)
		if !ok {
			return nil, ErrInvalidIntervalCoinCount
		}

		addedCoin, err := r.gen.Generate(rng, groupID, prevBoardResult)
		if err != nil {
			return nil, fmt.Errorf("generate new coin: %w", err)
		}

		coinCount += addedCoin
		res.SetBoardResult(prevBoardResult.Copy())

		if addedCoin > 0 {
			retrigger, err := handleReTrigger(rsCount)
			if err != nil {
				return nil, fmt.Errorf("handle retrigger: %w", err)
			}

			res.SetReTriggers(retrigger)
		}

		rsCount = r.applyRetrigger(rsCount, res)

		if rsCount == 0 {
			stepWin := r.coinManager.Collect(res.BoardResult(), eventer, winLimiter)

			res.SetStepWin(stepWin)
		}

		res.SetEvents(eventer.Get())
		res.SetTotalWin(winLimiter.TotalWin())
		res.SetMaxWinReached(winLimiter.MaxWinReached())

		allRes = append(allRes, res)

	}

	return allRes, nil

}

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
