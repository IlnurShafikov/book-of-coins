// Package free implements the free spins phase of the slot game.
// The free spins phase includes features such as expanding symbols,
// jackpot hunt coin collection, and retrigger mechanisms.
package free

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
	"github.com/IlnurShafikov/slot-engine/pkg/calculators"
	"github.com/IlnurShafikov/slot-engine/pkg/generator"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/fscounter"
)

// SpinPhase represents the free spins phase with all its components.
// It contains the game mechanics and configurations needed to execute free spins
// including reel generation, win calculation, symbol expansion, and jackpot features.
type SpinPhase struct {
	// gen generates reel symbols for each free spin
	gen *generator.Generator
	// calc performs line and scatter win calculations
	calc calculator.Calculators

	getterValues *getter.CoinValue

	expandMinCount map[int]int

	fSInitCounter fscounter.FSInitCounter

	// rsInitCounter determines initial re spins count
	rsInitCounter intervals.Intervals
}

// NewSpinPhase creates a new free spins phase with all required components.
// The returned SpinPhase is ready to execute free spin rounds.
//
// Parameters:
//   - gen: reel symbol generator for free spins
//   - calc: line and scatter win calculator
//   - fSInitCounter: counter that determines initial free spins grant amount
//
// Returns:
//   - pointer to a new SpinPhase instance
func NewSpinPhase(
	gen *generator.Generator,
	calc calculator.Calculators,
	getterValues *getter.CoinValue,
	minExpandCount map[int]int,
	fSInitCounter fscounter.FSInitCounter,
	rsInitCounter intervals.Intervals,
) *SpinPhase {
	return &SpinPhase{
		gen:            gen,
		calc:           calc,
		getterValues:   getterValues,
		expandMinCount: minExpandCount,
		fSInitCounter:  fSInitCounter,
		rsInitCounter:  rsInitCounter,
	}
}

func (f *SpinPhase) Spin(
	rng model.RNG,
	expandSymbolID int,
	winLimiter limiter.WinLimiter,
) (*result.Free, error) {
	res := result.NewFree(expandSymbolID)
	eventer := collectors.NewEvents()

	board, err := f.gen.Generate(rng, expandSymbolID)
	if err != nil {
		return nil, fmt.Errorf("generate board: %w", err)
	}

	coinRes, err := f.getterValues.GenerateValues(rng, board)
	if err != nil {
		return nil, fmt.Errorf("generate coin values: %w", err)
	}

	boardRes := &result.BoardResult{
		Board: board,
		Coins: coinRes,
	}

	res.SetBeforeBoard(boardRes)

	stepWin, err := f.processWin(res, eventer, winLimiter)
	if err != nil {
		return nil, fmt.Errorf("process line win: %w", err)
	}

	err = f.checkRetrigger(res)
	if err != nil {
		return nil, fmt.Errorf("check retrigger: %w", err)
	}

	res.SetStepWin(stepWin)
	res.SetTotalWin(winLimiter.TotalWin())
	res.SetMaxWinReached(winLimiter.MaxWinReached())
	res.SetEvents(eventer.Get())

	return res, nil
}

func (f *SpinPhase) checkRetrigger(res *result.Free) error {
	scatterCount := res.BeforeBoard().Board.CalculateSymbol(rules.ScatterID)
	fsCount := f.fSInitCounter.GetFsInitCount(scatterCount)

	safeFsValue, err := cast.SafeCast[uint32](fsCount)
	if err != nil {
		return fmt.Errorf("safe cast: %w", err)
	}

	res.SetReTriggers(safeFsValue)

	rsTrigger := f.checkExpandCoins(res)
	if !rsTrigger {
		return nil
	}

	coinCount := res.AfterBoard().Board.CalculateSymbol(rules.CoinID)

	rsCount, ok := f.rsInitCounter.GetIDByValue(coinCount)
	if !ok || rsCount == 0 {
		return nil
	}

	safeRsValue, err := cast.SafeCast[uint32](rsCount)
	if err != nil {
		return fmt.Errorf("safe cast: %w", err)
	}

	res.SetReSpinsInitial(safeRsValue)

	return nil
}

// processWin calculates all win types for a single free spin and collects jackpot coins.
// It evaluates line wins, scatter wins, expanded symbol wins, and coin jackpot results.
//
// Parameters:
//   - res: free spin result containing the board and expansion state
//   - eventer: event collector for recording win and expansion events
//   - winLimiter: limiter enforcing maximum win cap
//
// Returns:
//   - stepWin: total win amount for this spin step
//   - hourglassCount: count of hourglass symbols found (0 if expansion occurred)
//   - error if calculation fails
func (f *SpinPhase) processWin(
	res *result.Free,
	eventer *collectors.EventsCollector,
	winLimiter limiter.WinLimiter,
) (stepWin uint64, err error) {
	board := res.BeforeBoard().Board

	lineWins, err := f.calc.CalculateLines(board, winLimiter)
	if err != nil {
		return 0, fmt.Errorf("calculate win lines: %w", err)
	}

	scatterWins, err := f.calc.CalculateScatters(board, winLimiter)
	if err != nil {
		return 0, fmt.Errorf("calculate scatter win: %w", err)
	}

	if len(lineWins) > 0 || len(scatterWins) > 0 {
		eventer.Add(events.NewWinsEvent(lineWins, scatterWins))
	}

	err = f.checkExpandReels(res, eventer)
	if err != nil {
		return 0, fmt.Errorf("check expand reels: %w", err)
	}

	expandLineWin, err := f.calc.CalculateExpand(
		res.ExpandSymbolID(), res.ExpandReels(), winLimiter)
	if err != nil {
		return 0, fmt.Errorf("calculate expand: %w", err)
	}

	if len(expandLineWin) > 0 {
		eventer.Add(events.NewWinsEvent(expandLineWin, nil))
	}

	stepWin = stepWins(lineWins, scatterWins, expandLineWin)

	return stepWin, nil
}

func (f *SpinPhase) checkExpandReels(res *result.Free, eventer *collectors.EventsCollector) error {
	board := res.BeforeBoard().Board

	expandCount := board.CalculateSymbol(res.ExpandSymbolID())

	minCount, ok := f.expandMinCount[res.ExpandSymbolID()]
	if !ok {
		return ErrInvalidExpandSymbolID
	}

	if expandCount < minCount {
		return nil
	}

	expandPos := board.SymbolPositions(res.ExpandSymbolID())

	expandReel := make([]int, 0, len(expandPos))

	for _, pos := range expandPos {
		expandReel = append(expandReel, pos.ReelIndex)
	}

	res.SetExpandReels(expandReel)
	eventer.Add(events.NewExpandSymbolEvent(res.ExpandReels(), res.ExpandSymbolID()))

	return nil
}

func (f *SpinPhase) checkExpandCoins(res *result.Free) bool {
	if res.ExpandSymbolID() != rules.CoinID {
		return false
	}

	if len(res.ExpandReels()) < f.expandMinCount[res.ExpandSymbolID()] {
		return false
	}

	boardRes := res.BeforeBoard().Copy()
	rows := boardRes.Board.FirstReelSize()

	for _, reel := range res.ExpandReels() {
		for row := range rows {
			boardRes.Board[reel][row] = rules.CoinID
		}
	}

	newCoinRes := make(coins.CoinResult, 0, len(res.ExpandReels())*rows)

	for _, coin := range boardRes.Coins {
		for row := range rows {
			newCoin := coins.Coin{
				Position: model.NewPosition(coin.Position.ReelIndex, row),
				Value:    coin.Value,
				JpType:   coin.JpType,
				SymbolID: coin.SymbolID,
			}

			newCoinRes = append(newCoinRes, newCoin)
		}
	}

	boardRes.Coins = newCoinRes
	res.SetAfterBoard(boardRes)

	return true
}

// stepWins aggregates all win sources into a single step win total.
// It sums line wins, scatter wins, expansion wins, and coin wins.
//
// Parameters:
//   - lineWins: wins from matching paylines
//   - scatterWins: wins from scattered symbol matches
//   - expandLineWins: wins generated by expanded symbols
//
// Returns:
//   - total step win across all win types
func stepWins(
	lineWins calculators.LineWins,
	scatterWins calculators.PayAnyWhereWins,
	expandLineWins calculators.LineWins,
) uint64 {
	stepLineWins := stepWinInLines(lineWins)
	stepScatterWins := stepWinInScatters(scatterWins)
	stepExpandLineWins := stepWinInLines(expandLineWins)

	return stepLineWins + stepScatterWins + stepExpandLineWins
}

// stepWinInLines calculates the total win amount from all line wins.
//
// Parameters:
//   - wins: collection of line wins
//
// Returns:
//   - total payout from all wins
func stepWinInLines(wins calculators.LineWins) uint64 {
	if len(wins) == 0 {
		return 0
	}

	var totalWin uint64

	for _, win := range wins {
		totalWin += win.Payout
	}

	return totalWin
}

// stepWinInScatters calculates the total win amount from all scatter wins.
//
// Parameters:
//   - wins: collection of scatter wins
//
// Returns:
//   - total payout from all scatter wins
func stepWinInScatters(wins calculators.PayAnyWhereWins) uint64 {
	if len(wins) == 0 {
		return 0
	}

	var totalWin uint64

	for _, win := range wins {
		totalWin += win.Payout
	}

	return totalWin
}
