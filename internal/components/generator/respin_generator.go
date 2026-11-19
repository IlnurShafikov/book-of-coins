// Package generator contains the generator that drops new coins into the empty cells
// of the board during respins.
package generator

import (
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/components/getter"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/result"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	"github.com/IlnurShafikov/slot-engine/pkg/weight"
)

// ReSpinGenerator drops new coins into the empty cells of the board during respins.
type ReSpinGenerator struct {
	// getterValues draws the values of the new coins.
	getterValues *getter.CoinValue
	// needPlaceCoin holds, for every coin chance group, the weights of an empty cell staying
	// empty or becoming a coin.
	needPlaceCoin weight.Group
}

// NewReSpinGenerator creates a respin generator.
//
// Parameters:
//   - needPlaceCoin: weights of an empty cell staying empty or becoming a coin per group
//   - getterValues: draws the values of the new coins
//
// Returns:
//   - pointer to a new ReSpinGenerator
func NewReSpinGenerator(
	needPlaceCoin weight.Group,
	getterValues *getter.CoinValue,
) *ReSpinGenerator {
	return &ReSpinGenerator{
		needPlaceCoin: needPlaceCoin,
		getterValues:  getterValues,
	}
}

// Generate decides which empty cells of the board get a coin, places the coins and draws
// their values. The new coins are appended to boardResult.Coins.
//
// Parameters:
//   - rng: random number generator
//   - groupID: coin chance group chosen by the number of coins on the board
//   - boardResult: board of the respin phase (modified)
//
// Returns:
//   - number of coins added in this respin
//   - error if a coin cannot be placed or its value cannot be generated
func (r *ReSpinGenerator) Generate(
	rng model.RNG,
	groupID int,
	boardResult *result.BoardResult,
) (int, error) {
	placePos, err := r.placeCoinToBoard(rng, groupID, boardResult)
	if err != nil {
		return 0, fmt.Errorf("place coin to board: %w", err)
	}

	if len(placePos) == 0 {
		return 0, nil
	}

	newCoins, err := r.getterValues.GenerateValuesInPositions(rng, placePos)
	if err != nil {
		return 0, fmt.Errorf("generate coin values: %w", err)
	}

	boardResult.Coins = append(boardResult.Coins, newCoins...)

	return len(newCoins), nil
}

// placeCoinToBoard walks over the empty cells and independently turns each of them into
// a coin with the weight of the given group. All cells of one respin use the same group.
//
// Parameters:
//   - rng: random number generator
//   - groupID: coin chance group
//   - boardResult: board of the respin phase (modified: coins are set on the board)
//
// Returns:
//   - positions of the placed coins (nil if the board has no empty cells)
//   - error if the random value cannot be drawn
func (r *ReSpinGenerator) placeCoinToBoard(
	rng model.RNG,
	groupID int,
	boardResult *result.BoardResult,
) (model.Positions, error) {
	freePos := boardResult.Board.SymbolPositions(model.Void)

	if len(freePos) == 0 {
		return nil, nil
	}

	placePos := make(model.Positions, 0, len(freePos))

	for _, pos := range freePos {
		symbolID, err := r.needPlaceCoin.GetValue(rng, groupID)
		if err != nil {
			return nil, fmt.Errorf("get random value: %w", err)
		}

		if symbolID == model.Void {
			continue
		}

		boardResult.Board.Set(pos.ReelIndex, pos.RowIndex, symbolID)

		placePos = append(placePos, pos)
	}

	return placePos, nil
}
