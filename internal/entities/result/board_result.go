package result

import (
	"github.com/IlnurShafikov/book-of-coins/internal/entities/coins"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
)

// BoardResult represents the state of the game board with associated game elements.
// This structure combines the symbol and coin results.
type BoardResult struct {
	// Board contains the symbol matrix (reels x rows) for the current state.
	Board model.Board

	// Coins contains generated coin results associated with board positions.
	Coins coins.CoinResult
}

// NewBoardResult creates a new BoardResult with the specified board.
//
// Parameters:
//   - board: the game board containing symbol positions
//
// Returns:
//   - pointer to the newly created BoardResult
func NewBoardResult(board model.Board) *BoardResult {
	return &BoardResult{
		Board: board,
		Coins: nil,
	}
}

// Copy creates a copy of the BoardResult.
// Board, multipliers, unions, coins, and sticky frames are cloned,
// while PositionState reference is preserved.
//
// Returns:
//   - pointer to a new BoardResult with copied values
func (b *BoardResult) Copy() *BoardResult {
	return &BoardResult{
		Board: b.Board.Clone(),
		Coins: b.Coins.Copy(),
	}
}
