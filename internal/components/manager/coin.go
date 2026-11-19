// Package manager contains the component that collects the coins of a board and pays the
// coin win and the grand jackpot.
package manager

import (
	"github.com/IlnurShafikov/book-of-coins/internal/components/collectors"
	"github.com/IlnurShafikov/book-of-coins/internal/components/limiter"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/result"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/result/events"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/rules"
)

// Coin collects the coins of a board and pays the accumulated win, including the grand
// jackpot for a full board.
type Coin struct {
	// grandWin is the grand jackpot prize added to the coin sum.
	grandWin uint64
	// grandJpType is the jackpot type ID sent with the grand win event.
	grandJpType int
	// coinCountToGrand is the number of coins on the board that awards the grand jackpot.
	coinCountToGrand int
}

// NewCoin creates a coin win manager.
//
// Parameters:
//   - grandWin: grand jackpot prize
//   - grandJpType: jackpot type ID of the grand jackpot
//   - coinCountToGrand: number of coins on the board that awards the grand jackpot
//
// Returns:
//   - pointer to a new Coin
func NewCoin(
	grandWin uint64,
	grandJpType int,
	coinCountToGrand int,
) *Coin {
	return &Coin{
		grandWin:         grandWin,
		grandJpType:      grandJpType,
		coinCountToGrand: coinCountToGrand,
	}
}

// Collect sums the values of all coins of the board, adds the grand jackpot when the board
// holds coinCountToGrand coins, passes the sum through the win limiter and records the
// collect (and grand win) events.
//
// Parameters:
//   - res: board result with the coins to collect
//   - eventer: event collector for the client events
//   - winLimiter: limiter enforcing the max win cap
//
// Returns:
//   - win paid for the coins after the cap (0 if there are no coins)
func (c *Coin) Collect(
	res *result.BoardResult,
	eventer *collectors.EventsCollector,
	winLimiter limiter.WinLimiter,
) uint64 {
	if res.Coins == nil {
		return 0
	}

	var totalCoinSum uint64

	for _, coin := range res.Coins {
		totalCoinSum += coin.Value
	}

	if res.Board.CalculateSymbol(rules.CoinID) == c.coinCountToGrand {
		totalCoinSum += c.grandWin

		grandEvent := events.NewGrandWinEvent(c.grandWin, c.grandJpType)
		eventer.Add(grandEvent)
	}

	stepWin, _ := winLimiter.Add(totalCoinSum)

	collectEvent := events.NewCollectCoinEvent(res.Coins.Copy())
	eventer.Add(collectEvent)

	return stepWin
}
