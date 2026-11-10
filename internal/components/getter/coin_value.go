package getter

import (
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/entities/coins"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/rules"
	"github.com/IlnurShafikov/book-of-coins/internal/tools/cast"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	"github.com/IlnurShafikov/slot-engine/pkg/weight"
)

type CoinValue struct {
	coinValues *weight.Weights
	jpTypes    map[uint64]int
}

func NewCoinValue(
	coinValues *weight.Weights,
	jpTypes map[uint64]int,
) *CoinValue {
	return &CoinValue{
		coinValues: coinValues,
		jpTypes:    jpTypes,
	}
}

func (c *CoinValue) GenerateValues(
	rng model.RNG,
	board model.Board,
) (coins.CoinResult, error) {

	coinCount := board.CalculateSymbol(rules.CoinID)
	if coinCount == 0 {
		return nil, nil
	}

	coinPositions := board.SymbolPositions(rules.CoinID)

	resp, err := c.GenerateValuesInPositions(rng, coinPositions)
	if err != nil {
		return nil, fmt.Errorf("generate coin values: %w", err)
	}

	return resp, nil
}

func (c *CoinValue) GenerateValuesInPositions(
	rng model.RNG,
	positions model.Positions,
) (coins.CoinResult, error) {
	if len(positions) == 0 {
		return nil, nil
	}

	resp := make(coins.CoinResult, len(positions))

	for i, pos := range positions {
		value, err := c.coinValues.GetRandomValue(rng)
		if err != nil {
			return nil, fmt.Errorf("generate coin value: %w", err)
		}

		safeValue, err := cast.SafeCast[uint64](value)
		if err != nil {
			return nil, fmt.Errorf("safe int to uint64: %w", err)
		}

		jpType, _ := c.jpTypes[safeValue]

		resp[i] = coins.Coin{
			Position: pos,
			Value:    safeValue,
			JpType:   jpType,
			SymbolID: rules.CoinID,
		}
	}

	return resp, nil
}
