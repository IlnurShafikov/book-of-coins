package calculator

import (
	"fmt"
	"slices"

	"github.com/IlnurShafikov/book-of-coins/internal/components/limiter"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/rules"
	"github.com/IlnurShafikov/slot-engine/pkg/calculators"
	"github.com/IlnurShafikov/slot-engine/pkg/lines"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	"github.com/IlnurShafikov/slot-engine/pkg/paytable"
)

// ExpandCalculator calculates line wins for expanding symbols.
// An expanding symbol fills its entire reel; wins are awarded on every
// configured line based on how many reels contain the symbol.
type ExpandCalculator struct {
	// payTable holds the line payouts of the symbols.
	payTable *paytable.PayTable
	// lines holds the payline definitions the expansion is paid on.
	lines lines.Lines
}

// NewExpandCalculator creates an ExpandCalculator with the given paytable and line definitions.
func NewExpandCalculator(payTable *paytable.PayTable, winLines lines.Lines) *ExpandCalculator {
	return &ExpandCalculator{payTable: payTable, lines: winLines}
}

// CalculateExpand returns line wins for the expanding symbol on the board.
// expandReels contains 0-based indices of reels where the symbol is present.
// Returns nil without error when the reel count falls below the paytable minimum.
func (e *ExpandCalculator) CalculateExpand(
	symbolID model.Symbol,
	expandReels []int,
	winLimiter limiter.WinLimiter,
) (calculators.LineWins, error) {
	if len(expandReels) == 0 || symbolID == rules.CoinID {
		return nil, nil
	}

	payout, err := e.payTable.GetPayout(symbolID, len(expandReels))
	if err != nil {
		return nil, fmt.Errorf("get payout: %w", err)
	}

	if payout == 0 {
		return nil, nil
	}

	wins := make(calculators.LineWins, 0, len(e.lines))

	for _, line := range e.lines {
		limitedPayout, maxReached := winLimiter.Add(payout)

		wins = append(wins, calculators.LineWin{
			Symbol:    symbolID,
			Number:    line.Number,
			Payout:    limitedPayout,
			Positions: positionsOnExpandedReels(line, expandReels),
		})

		if maxReached {
			break
		}
	}

	return wins, nil
}

// positionsOnExpandedReels returns line positions that fall on expanded reels.
func positionsOnExpandedReels(line lines.Line, expanded []int) model.Positions {
	positions := make(model.Positions, 0, len(expanded))

	for _, pos := range line.Positions {
		if slices.Contains(expanded, pos.ReelIndex) {
			positions = append(positions, pos)
		}
	}

	return positions
}
