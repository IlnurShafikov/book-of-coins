package datatransfer

import (
	"github.com/IlnurShafikov/book-of-coins/internal/entities/coins"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/result"
	"github.com/IlnurShafikov/book-of-coins/pkg/datatransfer/dto/common"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
)

// TransformSymbols converts the internal game board representation into client DTO symbols.
// This function should be implemented to transform your game's board structure into
// an array of Symbol DTOs that will be sent to the client.
//
// Parameters:
//   - gameBoard: the internal game board (type depends on your implementation)
//   - betValue: the player's bet value (may be needed for property calculations)
//
// Returns:
//   - slice of Symbol DTOs representing the board
//   - error if transformation fails
func TransformSymbols(gameBoard *result.BoardResult, betValue uint64) ([]*common.Symbol, error) {
	if gameBoard == nil {
		return nil, ErrGameBoardIsNil
	}

	resp := transformBoardToSymbols(gameBoard)

	updateSymbolsWithCoins(gameBoard.Coins, resp, betValue)

	return resp, nil
}

// transformBoardToSymbols converts the board result into a flat array of Symbol DTOs.
// Each symbol includes its ID, position on the board, and optional properties.
//
// Parameters:
//   - boardResult: the board result containing the symbol layout
//
// Returns:
//   - slice of Symbol DTOs in row-major order (reel by reel, then row by row)
func transformBoardToSymbols(gameBoard *result.BoardResult) []*common.Symbol {
	reelsCount := gameBoard.Board.ReelCount()
	rowsCount := gameBoard.Board.FirstReelSize()

	resp := make([]*common.Symbol, 0, reelsCount*rowsCount)

	for reel := range reelsCount {
		for row := range rowsCount {
			resp = append(resp, &common.Symbol{
				SymbolID: gameBoard.Board[reel][row],
				Position: model.NewPosition(reel, row),
				Property: nil,
			})
		}
	}

	return resp
}

// updateSymbolsWithCoins attaches coin properties (value and jackpot type)
// to symbols at matching positions.
// Coins are matched by board position, and their values are scaled by bet amount.
//
// Parameters:
//   - coinResult: collection of coins to attach to symbols
//   - symbols: slice of Symbol DTOs to update (modified in-place)
//   - betValue: the player's bet amount to scale coin values
func updateSymbolsWithCoins(
	coinResult coins.CoinResult,
	symbols []*common.Symbol,
	betValue uint64,
) {
	if len(coinResult) == 0 {
		return
	}

	coinByPosition := make(map[model.Position]coins.Coin, len(coinResult))

	for _, coin := range coinResult {
		coinByPosition[coin.Position] = coin
	}

	for i, symbol := range symbols {
		coin, ok := coinByPosition[symbol.Position]
		if ok {
			symbols[i].Property = &common.SymbolProperties{
				Value:  coin.Value * betValue,
				JpType: coin.JpType,
			}
		}
	}
}
