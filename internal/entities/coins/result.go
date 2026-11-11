// Package coins defines data structures for coin results on the game board.
package coins

import "github.com/IlnurShafikov/slot-engine/pkg/model"

// CoinResult is a collection of Coin items representing coins on the game board.
type CoinResult []Coin

// Coin represents a coin symbol on the game board with its associated value and properties.
type Coin struct {
	// Position is the location of the coin on the game board
	Position model.Position `json:"position"`

	// Value is the monetary value of the coin
	Value uint64 `json:"value"`

	// JpType is the jackpot type identifier (0 if not a jackpot coin)
	JpType int `json:"jpType,omitempty"`

	// SymbolID is the symbol identifier (CoinID or GoldCoinID)
	SymbolID int `json:"symbol"`
}

// Copy creates a shallow copy of the Coin.
// Since all fields are value types, a shallow copy is sufficient.
//
// Returns:
//   - a new Coin with copied field values
func (c Coin) Copy() Coin {
	return Coin{
		Position: c.Position,
		Value:    c.Value,
		JpType:   c.JpType,
		SymbolID: c.SymbolID,
	}
}

// Copy creates a deep copy of the CoinResult collection.
//
// Returns:
//   - a new CoinResult slice with copied Coin elements, or nil if the original is nil
func (cr CoinResult) Copy() CoinResult {
	if cr == nil {
		return nil
	}

	res := make(CoinResult, len(cr))
	for i := range cr {
		res[i] = cr[i].Copy()
	}

	return res
}
