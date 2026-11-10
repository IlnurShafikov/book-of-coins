// Package calculator provides win calculation functionality for game lines.
package calculator

import (
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/components/limiter"
	"github.com/IlnurShafikov/slot-engine/pkg/calculators"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
)

// Calculators defines the interface for calculating line-based and scatter wins on a game board.
type Calculators interface {
	// CalculateLines processes the board and returns all winning lines
	CalculateLines(
		board model.Board, winLimiter limiter.WinLimiter,
	) (calculators.LineWins, error)

	// CalculateScatters processes the board and returns all scatter wins
	CalculateScatters(
		board model.Board, winLimiter limiter.WinLimiter,
	) (calculators.PayAnyWhereWins, error)

	// CalculateExpand returns line wins for an expanding symbol across the given reels
	CalculateExpand(
		symbolID model.Symbol,
		expandReels []int,
		winLimiter limiter.WinLimiter,
	) (calculators.LineWins, error)
}

// Calculator implements the Calculators interface using configured pay lines and scatter paytable.
type Calculator struct {
	// lines contains the pay line definitions and paytable
	lines *calculators.Lines
	// payAnyWhere contains the scatter paytable for scatter wins
	payAnyWhere *calculators.PayAnyWhere
	// expandCalculator handles win calculation for expanding symbols
	expandCalculator *ExpandCalculator
}

// NewCalculators creates a new Calculator with the specified pay lines,
// scatter paytable, and expand calculator.
//
// Parameters:
//   - lines: configured pay lines and paytable for line wins
//   - payAnyWhere: configured scatter paytable for scatter wins
//   - expandCalculator: calculator for expanding symbol line wins
//
// Returns:
//   - pointer to a new Calculator
func NewCalculators(
	lines *calculators.Lines,
	payAnyWhere *calculators.PayAnyWhere,
	expandCalculator *ExpandCalculator,
) *Calculator {
	return &Calculator{
		lines:            lines,
		payAnyWhere:      payAnyWhere,
		expandCalculator: expandCalculator,
	}
}

// CalculateLines processes the game board and returns all winning line combinations.
// It applies the win limiter to ensure wins don't exceed maximum limits.
//
// Parameters:
//   - board: the game board to calculate wins on
//   - winLimiter: limiter to cap maximum wins
//
// Returns:
//   - calculators.LineWins containing all winning lines
//   - error if calculation fails
func (c *Calculator) CalculateLines(
	board model.Board,
	winLimiter limiter.WinLimiter,
) (calculators.LineWins, error) {
	lineWins, err := c.lines.CalculateWithStrategies(
		board, calculators.WithLinesWinLimiter(winLimiter),
	)
	if err != nil {
		return nil, fmt.Errorf("calculate lines: %w", err)
	}

	return lineWins, nil
}

// CalculateExpand returns line wins for the expanding symbol across the given reels.
// Delegates to ExpandCalculator, applying the win limiter to cap maximum wins.
//
// Parameters:
//   - symbolID: the symbol that has expanded to fill its reels
//   - expandReels: 0-based indices of reels where the symbol is present
//   - winLimiter: limiter to cap maximum wins
//
// Returns:
//   - calculators.LineWins for the expanded symbol, nil if count is below paytable minimum
//   - error if payout lookup fails
func (c *Calculator) CalculateExpand(
	symbolID model.Symbol,
	expandReels []int,
	winLimiter limiter.WinLimiter,
) (calculators.LineWins, error) {
	wins, err := c.expandCalculator.CalculateExpand(symbolID, expandReels, winLimiter)
	if err != nil {
		return nil, fmt.Errorf("calculate expand: %w", err)
	}

	return wins, nil
}

// CalculateScatters processes the game board and returns all scatter wins.
// It applies the win limiter to ensure wins don't exceed maximum limits.
//
// Parameters:
//   - board: the game board to calculate scatter wins on
//   - winLimiter: limiter to cap maximum wins
//
// Returns:
//   - calculators.PayAnyWhereWins containing all scatter wins
//   - error if calculation fails
func (c *Calculator) CalculateScatters(
	board model.Board, winLimiter limiter.WinLimiter,
) (calculators.PayAnyWhereWins, error) {
	wins, err := c.payAnyWhere.CalculateWithStrategies(
		board, calculators.WithPayAnyWhereWinLimiter(winLimiter),
	)
	if err != nil {
		return nil, fmt.Errorf("calculate payAnyWhere: %w", err)
	}

	return wins, nil
}
