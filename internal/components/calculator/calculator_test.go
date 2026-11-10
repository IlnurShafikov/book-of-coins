package calculator_test

import (
	"math"
	"testing"

	"github.com/IlnurShafikov/book-of-coins/internal/components/calculator"
	"github.com/IlnurShafikov/slot-engine/pkg/calculators"
	"github.com/IlnurShafikov/slot-engine/pkg/lines"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	"github.com/IlnurShafikov/slot-engine/pkg/paytable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	symPaying  = model.Symbol(1)
	symNoPay   = model.Symbol(2)
	symScatter = model.Symbol(5)
)

var (
	lineWinBoard = model.Board{
		{symPaying, symNoPay, symNoPay},
		{symPaying, symNoPay, symNoPay},
		{symPaying, symNoPay, symNoPay},
		{symPaying, symNoPay, symNoPay},
		{symPaying, symNoPay, symNoPay},
	}

	lineNoWinBoard = model.Board{
		{symNoPay, symNoPay, symNoPay},
		{symNoPay, symNoPay, symNoPay},
		{symNoPay, symNoPay, symNoPay},
		{symNoPay, symNoPay, symNoPay},
		{symNoPay, symNoPay, symNoPay},
	}

	scatterWinBoard = model.Board{
		{symNoPay, symScatter, symNoPay},
		{symNoPay, symScatter, symNoPay},
		{symNoPay, symScatter, symNoPay},
		{symNoPay, symNoPay, symNoPay},
		{symNoPay, symNoPay, symNoPay},
	}

	scatterNoWinBoard = model.Board{
		{symNoPay, symNoPay, symNoPay},
		{symNoPay, symNoPay, symNoPay},
		{symNoPay, symNoPay, symNoPay},
		{symNoPay, symNoPay, symNoPay},
		{symNoPay, symNoPay, symNoPay},
	}
)

func newMainCalculator(t *testing.T) *calculator.Calculator {
	t.Helper()

	linePT, err := paytable.NewPayTable(
		paytable.SymbolsPayouts{
			symPaying: paytable.SymbolPayouts{0, 100, 200, 500},
			symNoPay:  paytable.SymbolPayouts{0, 0, 0, 0},
		},
		paytable.CountRanges{
			paytable.NewRange(2, 2),
			paytable.NewRange(3, 3),
			paytable.NewRange(4, 4),
			paytable.NewRange(5, 5),
		},
	)
	require.NoError(t, err)

	ls := lines.Lines{
		lines.NewLine(1, model.Positions{
			model.NewPosition(0, 0), model.NewPosition(1, 0), model.NewPosition(2, 0),
			model.NewPosition(3, 0), model.NewPosition(4, 0),
		}),
		lines.NewLine(2, model.Positions{
			model.NewPosition(0, 1), model.NewPosition(1, 1), model.NewPosition(2, 1),
			model.NewPosition(3, 1), model.NewPosition(4, 1),
		}),
		lines.NewLine(3, model.Positions{
			model.NewPosition(0, 2), model.NewPosition(1, 2), model.NewPosition(2, 2),
			model.NewPosition(3, 2), model.NewPosition(4, 2),
		}),
	}

	scatterPT, err := paytable.NewPayTable(
		paytable.SymbolsPayouts{
			symScatter: paytable.SymbolPayouts{50, 200, 500},
		},
		paytable.CountRanges{
			paytable.NewRange(3, 3),
			paytable.NewRange(4, 4),
			paytable.NewRange(5, 5),
		},
	)
	require.NoError(t, err)

	return calculator.NewCalculators(
		calculators.NewLines(linePT, ls, nil, nil),
		calculators.NewPayAnyWhere(scatterPT, model.Collection{symScatter}),
		nil,
	)
}

// newCalculatorWithInvalidLine creates a Calculator whose line 1 has an
// out-of-bounds reel index, causing CalculateLines to return an error.
func newCalculatorWithInvalidLine(t *testing.T) *calculator.Calculator {
	t.Helper()

	linePT, err := paytable.NewPayTable(
		paytable.SymbolsPayouts{
			symPaying: paytable.SymbolPayouts{0, 100, 200, 500},
			symNoPay:  paytable.SymbolPayouts{0, 0, 0, 0},
		},
		paytable.CountRanges{
			paytable.NewRange(2, 2),
			paytable.NewRange(3, 3),
			paytable.NewRange(4, 4),
			paytable.NewRange(5, 5),
		},
	)
	require.NoError(t, err)

	// reel 5 does not exist on a 5-reel board → GetSymbol returns error.
	invalidLine := lines.NewLine(1, model.Positions{
		model.NewPosition(5, 0), model.NewPosition(1, 0), model.NewPosition(2, 0),
		model.NewPosition(3, 0), model.NewPosition(4, 0),
	})

	scatterPT, err := paytable.NewPayTable(
		paytable.SymbolsPayouts{
			symScatter: paytable.SymbolPayouts{50, 200, 500},
		},
		paytable.CountRanges{
			paytable.NewRange(3, 3),
			paytable.NewRange(4, 4),
			paytable.NewRange(5, 5),
		},
	)
	require.NoError(t, err)

	return calculator.NewCalculators(
		calculators.NewLines(linePT, lines.Lines{invalidLine}, nil, nil),
		calculators.NewPayAnyWhere(scatterPT, model.Collection{symScatter}),
		nil,
	)
}

func TestNewCalculators(t *testing.T) {
	t.Parallel()

	require.NotNil(t, newMainCalculator(t))
}

func TestCalculator_CalculateLines(t *testing.T) {
	t.Parallel()

	calc := newMainCalculator(t)

	tests := []struct {
		name       string
		board      model.Board
		wantWins   int
		wantSym    model.Symbol
		wantPayout uint64
	}{
		{
			name:       "5-of-kind on line 1",
			board:      lineWinBoard,
			wantWins:   1,
			wantSym:    symPaying,
			wantPayout: 500,
		},
		{
			name:     "all zero-payout symbol — no wins",
			board:    lineNoWinBoard,
			wantWins: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			wins, err := calc.CalculateLines(tc.board, newLimiter(t, math.MaxUint64))
			require.NoError(t, err)
			require.Len(t, wins, tc.wantWins)

			for _, w := range wins {
				assert.Equal(t, tc.wantSym, w.Symbol)
				assert.Equal(t, tc.wantPayout, w.Payout)
			}
		})
	}
}

func TestCalculator_CalculateLines_Error(t *testing.T) {
	t.Parallel()

	_, err := newCalculatorWithInvalidLine(t).CalculateLines(lineWinBoard, newLimiter(t, math.MaxUint64))
	require.Error(t, err)
}

func TestCalculator_CalculateScatters(t *testing.T) {
	t.Parallel()

	calc := newMainCalculator(t)

	tests := []struct {
		name       string
		board      model.Board
		wantWins   int
		wantPayout uint64
	}{
		{
			name:       "3 scatters anywhere — payout 50",
			board:      scatterWinBoard,
			wantWins:   1,
			wantPayout: 50,
		},
		{
			name:     "scatter absent — no win",
			board:    scatterNoWinBoard,
			wantWins: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			wins, err := calc.CalculateScatters(tc.board, newLimiter(t, math.MaxUint64))
			require.NoError(t, err)
			require.Len(t, wins, tc.wantWins)

			for _, w := range wins {
				assert.Equal(t, tc.wantPayout, w.Payout)
			}
		})
	}
}

func TestCalculator_CalculateScatters_Error(t *testing.T) {
	t.Parallel()

	// 6 scatters on a 5x3 board; paytable covers only 3-5 → GetPayout(scatter, 6) errors.
	board := model.Board{
		{symScatter, symScatter, symNoPay},
		{symScatter, symScatter, symNoPay},
		{symScatter, symScatter, symNoPay},
		{symNoPay, symNoPay, symNoPay},
		{symNoPay, symNoPay, symNoPay},
	}

	_, err := newMainCalculator(t).CalculateScatters(board, newLimiter(t, math.MaxUint64))
	require.Error(t, err)
}
