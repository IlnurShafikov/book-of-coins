package manager

import (
	"github.com/IlnurShafikov/book-of-coins/internal/components/collectors"
	"github.com/IlnurShafikov/book-of-coins/internal/components/limiter"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/result"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/result/events"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/rules"
)

type Coin struct {
	grandWin         uint64
	grandJpType      int
	coinCountToGrand int
}

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
