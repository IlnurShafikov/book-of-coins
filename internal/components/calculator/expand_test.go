package calculator_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/IlnurShafikov/book-of-coins/internal/components/calculator"
	"github.com/IlnurShafikov/book-of-coins/internal/components/limiter"
	"github.com/IlnurShafikov/slot-engine/pkg/lines"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	"github.com/IlnurShafikov/slot-engine/pkg/paytable"
)

const expandSym = model.Symbol(7)

func newTestCalculator(t *testing.T) *calculator.ExpandCalculator {
	t.Helper()

	pt, err := paytable.NewPayTable(
		paytable.SymbolsPayouts{
			expandSym: paytable.SymbolPayouts{100, 200, 500},
		},
		paytable.CountRanges{
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

	return calculator.NewExpandCalculator(pt, ls)
}

func newLimiter(t *testing.T, maxWin uint64) limiter.WinLimiter {
	t.Helper()

	lim, err := limiter.New(maxWin)
	require.NoError(t, err)

	return lim
}

func TestExpandCalculator_Calculate(t *testing.T) {
	t.Parallel()

	calc := newTestCalculator(t)

	tests := []struct {
		name        string
		expandReels []int
		wantWins    int
		wantPayout  uint64
		wantPosCnt  int
	}{
		{
			name:        "3 reels — pays 3-of-a-kind on all lines",
			expandReels: []int{0, 2, 4},
			wantWins:    3,
			wantPayout:  100,
			wantPosCnt:  3,
		},
		{
			name:        "4 reels — pays 4-of-a-kind on all lines",
			expandReels: []int{0, 1, 2, 3},
			wantWins:    3,
			wantPayout:  200,
			wantPosCnt:  4,
		},
		{
			name:        "5 reels — pays 5-of-a-kind on all lines",
			expandReels: []int{0, 1, 2, 3, 4},
			wantWins:    3,
			wantPayout:  500,
			wantPosCnt:  5,
		},
		{
			name:        "empty expandReels — no win",
			expandReels: nil,
			wantWins:    0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			wins, err := calc.CalculateExpand(expandSym, tc.expandReels, newLimiter(t, math.MaxUint64))
			require.NoError(t, err)
			require.Len(t, wins, tc.wantWins)

			for _, w := range wins {
				assert.Equal(t, expandSym, w.Symbol)
				assert.Equal(t, tc.wantPayout, w.Payout)
				assert.Len(t, w.Positions, tc.wantPosCnt)
			}
		})
	}
}

func TestExpandCalculator_Calculate_Error(t *testing.T) {
	t.Parallel()

	calc := newTestCalculator(t)

	_, err := calc.CalculateExpand(12, []int{0, 2}, newLimiter(t, math.MaxUint64))
	require.Error(t, err)
}

func TestExpandCalculator_Calculate_LineNumbers(t *testing.T) {
	t.Parallel()

	calc := newTestCalculator(t)

	wins, err := calc.CalculateExpand(expandSym, []int{0, 2, 4}, newLimiter(t, math.MaxUint64))
	require.NoError(t, err)
	require.Len(t, wins, 3)

	for i, w := range wins {
		assert.Equal(t, i+1, w.Number)
	}
}

func TestExpandCalculator_Calculate_WinLimiter(t *testing.T) {
	t.Parallel()

	calc := newTestCalculator(t)

	wins, err := calc.CalculateExpand(expandSym, []int{0, 2, 4}, newLimiter(t, 150))
	require.NoError(t, err)
	require.Len(t, wins, 2)

	assert.Equal(t, uint64(100), wins[0].Payout)
	assert.Equal(t, uint64(50), wins[1].Payout)
}
