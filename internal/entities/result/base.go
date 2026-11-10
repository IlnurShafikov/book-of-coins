// Package result contains data structures for representing game results
// across different game phases (base, free spins, respins).
package result

import (
	"github.com/IlnurShafikov/book-of-coins/internal/entities/result/events"
)

// Base represents the result of a base game spin.
// It contains all information about the spin outcome including
// events, wins, and triggered features.
type Base struct {
	// boardResult contains the board state and reel multiplier result
	boardResult *BoardResult

	// events store all events that occurred during this spin
	events []events.Event

	// ReSpins contains results of respins triggered by this base spin
	ReSpins []*ReSpin

	// totalWin is the total win amount from this base spin (excluding respins/free spins)
	totalWin uint64

	// maxWinReached indicates if the maximum win limit was reached
	maxWinReached bool

	// expandSymbolID is the ID of the symbol chosen to expand during a free spins trigger.
	expandSymbolID int

	// freeSpinsInitials is the number of free spins initially triggered
	freeSpinsInitials uint32

	// reSpinInitials is the number of respins initially triggered
	reSpinInitials uint32
}

// NewBase creates a new Base result with default values.
//
// Returns:
//   - pointer to a new Base instance
func NewBase() *Base {
	return &Base{
		boardResult:       nil,
		events:            nil,
		ReSpins:           nil,
		totalWin:          0,
		maxWinReached:     false,
		expandSymbolID:    0,
		freeSpinsInitials: 0,
		reSpinInitials:    0,
	}
}

// SetBoardResult sets the board result containing the board state and reel multipliers.
//
// Parameters:
//   - board: pointer to BoardResult with the board state and calculation results
func (b *Base) SetBoardResult(board *BoardResult) {
	b.boardResult = board
}

// SetReSpinsInitial sets the initial number of respins triggered.
//
// Parameters:
//   - count: number of respins triggered
func (b *Base) SetReSpinsInitial(count uint32) {
	b.reSpinInitials = count
}

// SetEvents sets the events that occurred during this spin.
//
// Parameters:
//   - e: slice of events to set
func (b *Base) SetEvents(e []events.Event) {
	b.events = e
}

// SetFreeSpinsInitialsCount sets the number of free spins triggered.
//
// Parameters:
//   - count: number of free spins triggered
func (b *Base) SetFreeSpinsInitialsCount(count uint32) {
	b.freeSpinsInitials = count
}

// SetExpandSymbolID stores the ID of the symbol selected to expand during a free spins trigger.
//
// Parameters:
//   - symbolID: symbol identifier chosen as the expanding symbol for the free phase
func (b *Base) SetExpandSymbolID(symbolID int) {
	b.expandSymbolID = symbolID
}

// SetTotalWin sets the total win amount for this spin.
//
// Parameters:
//   - win: total win amount
func (b *Base) SetTotalWin(win uint64) {
	b.totalWin = win
}

// SetMaxWinReached sets whether the maximum win limit was reached.
//
// Parameters:
//   - reached: true if max win was reached, false otherwise
func (b *Base) SetMaxWinReached(reached bool) {
	b.maxWinReached = reached
}

// BoardResult returns the board result containing the board state and reel multipliers.
//
// Returns:
//   - pointer to BoardResult with the board state and calculation results
func (b *Base) BoardResult() *BoardResult {
	return b.boardResult
}

// Events returns all events that occurred during this spin.
//
// Returns:
//   - slice of events
func (b *Base) Events() []events.Event {
	return b.events
}

// ReSpinInitials returns the number of respins initially triggered.
//
// Returns:
//   - number of respins
func (b *Base) ReSpinInitials() uint32 {
	return b.reSpinInitials
}

// FreeSpinsInitialsCount returns the number of free spins triggered.
//
// Returns:
//   - number of free spins
func (b *Base) FreeSpinsInitialsCount() uint32 {
	return b.freeSpinsInitials
}

// ExpandSymbolID returns the ID of the symbol chosen to expand during a free spins trigger.
//
// Returns:
//   - symbol ID of the expanding symbol, or 0 if free spins were not triggered
func (b *Base) ExpandSymbolID() int {
	return b.expandSymbolID
}

// TotalWin returns the total win amount for this spin.
//
// Returns:
//   - total win amount
func (b *Base) TotalWin() uint64 {
	return b.totalWin
}

// MaxWinReached returns whether the maximum win limit was reached.
//
// Returns:
//   - true if max win was reached, false otherwise
func (b *Base) MaxWinReached() bool {
	return b.maxWinReached
}
