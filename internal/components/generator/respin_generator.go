package generator

import (
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/components/getter"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/result"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	"github.com/IlnurShafikov/slot-engine/pkg/weight"
)

type ReSpinGenerator struct {
	needPlaceCoin weight.Group
	getterValues  *getter.CoinValue
}

func NewReSpinGenerator(
	needPlaceCoin weight.Group,
	getterValues *getter.CoinValue,
) *ReSpinGenerator {
	return &ReSpinGenerator{
		needPlaceCoin: needPlaceCoin,
		getterValues:  getterValues,
	}
}

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
