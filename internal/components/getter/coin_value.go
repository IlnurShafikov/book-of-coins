// Package getter draws the values of the coins that land on the board.
package getter

import (
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/entities/coins"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/rules"
	"github.com/IlnurShafikov/book-of-coins/internal/tools/cast"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	"github.com/IlnurShafikov/slot-engine/pkg/weight"
)

// CoinValue draws the values of coins and attaches the jackpot type to jackpot values.
type CoinValue struct {
	// coinValues holds the weights of the coin values of the phase.
	coinValues *weight.Weights
	// jpTypes maps a jackpot coin value to its jackpot type ID.
	jpTypes map[uint64]int
}

// NewCoinValue creates a coin value getter.
//
// Parameters:
//   - coinValues: weights of the coin values of the phase
//   - jpTypes: map from a jackpot value to its jackpot type ID
//
// Returns:
//   - pointer to a new CoinValue
func NewCoinValue(
	coinValues *weight.Weights,
	jpTypes map[uint64]int,
) *CoinValue {
	return &CoinValue{
		coinValues: coinValues,
		jpTypes:    jpTypes,
	}
}

// GenerateValues draws a value for every coin symbol that lies on the board.
//
// Parameters:
//   - rng: random number generator
//   - board: board with the landed coins
//
// Returns:
//   - coins with their positions and values (nil if there are no coins)
//   - error if a value cannot be generated
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

// GenerateValuesInPositions draws an independent value for a coin at each position.
// A value that equals a jackpot value also gets the jackpot type.
//
// Parameters:
//   - rng: random number generator
//   - positions: positions of the coins
//
// Returns:
//   - coins with their positions and values (nil if there are no positions)
//   - error if a value cannot be generated or converted to uint64
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
