package result

import "github.com/IlnurShafikov/book-of-coins/internal/entities/result/events"

// Free represents the result of a single free spin.
// It contains all information about the free spin outcome, including events, wins, and retriggers.
type Free struct {
	beforeBoard *BoardResult
	afterBoard  *BoardResult
	// events accumulates all events that occurred during this free spin.
	events []events.Event
	// ReSpins contains results of respins triggered during this free spin
	ReSpins []*ReSpin
	// expandReels contains 0-based indices of reels where the expand symbol was placed.
	expandReels []int
	// fsTotalWin is the cumulative win amount up to and including this free spin.
	fsTotalWin uint64
	// stepWin is the win amount for this specific free spin.
	stepWin uint64
	// expandSymbolID is the ID of the symbol chosen to expand during a free spins trigger.
	expandSymbolID int
	// reTriggers is the number of additional free spins triggered by this spin.
	reTriggers uint32
	// reSpinInitial is the number of respins initially triggered by this free spin
	reSpinInitial uint32
	// maxWinReached indicates if the maximum win limit was reached.
	maxWinReached bool
}

// NewFree creates a new Free spin result initialized with the provided board and expansion state.
//
// Parameters:
//   - expandSymbolID: ID of the symbol chosen to expand during free spins trigger
//
// Returns:
//   - pointer to a new Free instance with counters and flags initialized to zero/false
func NewFree(expandSymbolID int) *Free {
	return &Free{
		beforeBoard:    nil,
		afterBoard:     nil,
		expandReels:    nil,
		expandSymbolID: expandSymbolID,
		events:         nil,
		ReSpins:        nil,
		fsTotalWin:     0,
		stepWin:        0,
		reTriggers:     0,
		reSpinInitial:  0,
		maxWinReached:  false,
	}
}

// SetBeforeBoard sets the board state before transformations in this free spin.
//
// Parameters:
//   - boardRes: board result before any symbol transformations
func (f *Free) SetBeforeBoard(boardRes *BoardResult) {
	f.beforeBoard = boardRes
}

// SetAfterBoard sets the board state after transformations in this free spin.
//
// Parameters:
//   - boardRes: board result after all symbol transformations
func (f *Free) SetAfterBoard(boardRes *BoardResult) {
	f.afterBoard = boardRes
}

// BeforeBoard returns the board state before transformations.
//
// Returns:
//   - pointer to BoardResult with the pre-transformation state
func (f *Free) BeforeBoard() *BoardResult {
	return f.beforeBoard
}

// AfterBoard returns the board state after transformations.
//
// Returns:
//   - pointer to BoardResult with the post-transformation state
func (f *Free) AfterBoard() *BoardResult {
	return f.afterBoard
}

// ExpandSymbolID returns the ID of the symbol chosen to expand during a free spins trigger.
//
// Returns:
//   - symbol ID of the expanding symbol, or 0 if free spins were not triggered
func (f *Free) ExpandSymbolID() int {
	return f.expandSymbolID
}

// ExpandReels returns the 0-based indices of reels where the expand symbol was placed.
//
// Returns:
//   - slice of reel indices with expansion placement
func (f *Free) ExpandReels() []int {
	return f.expandReels
}

func (f *Free) SetExpandReels(expandReels []int) {
	f.expandReels = expandReels
}

// SetEvents sets the events that occurred during this free spin.
//
// Parameters:
//   - e: slice of events to set
func (f *Free) SetEvents(e []events.Event) {
	f.events = e
}

// SetReTriggers sets the number of additional free spins triggered.
//
// Parameters:
//   - count: number of additional free spins
func (f *Free) SetReTriggers(count uint32) {
	f.reTriggers = count
}

// SetTotalWin sets the cumulative win amount.
//
// Parameters:
//   - win: cumulative win amount
func (f *Free) SetTotalWin(win uint64) {
	f.fsTotalWin = win
}

// SetMaxWinReached sets whether the maximum win limit was reached.
//
// Parameters:
//   - reached: true if max win was reached, false otherwise
func (f *Free) SetMaxWinReached(reached bool) {
	f.maxWinReached = reached
}

// SetStepWin sets the win amount for this specific free spin.
//
// Parameters:
//   - stepWin: win amount for this spin
func (f *Free) SetStepWin(stepWin uint64) {
	f.stepWin = stepWin
}

// SetReSpinsInitial sets the number of respins initially triggered.
//
// Parameters:
//   - count: number of respins triggered
func (f *Free) SetReSpinsInitial(count uint32) {
	f.reSpinInitial = count
}

// ReSpinInitials returns the number of respins initially triggered.
//
// Returns:
//   - number of respins
func (f *Free) ReSpinInitials() uint32 {
	return f.reSpinInitial
}

// Events returns all events that occurred during this free spin.
//
// Returns:
//   - slice of events
func (f *Free) Events() []events.Event {
	return f.events
}

// ReTriggers returns the number of additional free spins triggered.
//
// Returns:
//   - number of retriggers
func (f *Free) ReTriggers() uint32 {
	return f.reTriggers
}

// TotalWin returns the cumulative win amount.
//
// Returns:
//   - cumulative win amount
func (f *Free) TotalWin() uint64 {
	return f.fsTotalWin
}

// MaxWinReached returns whether the maximum win limit was reached.
//
// Returns:
//   - true if max win was reached, false otherwise
func (f *Free) MaxWinReached() bool {
	return f.maxWinReached
}

// StepWin returns the win amount for this specific free spin.
//
// Returns:
//   - win amount for this spin
func (f *Free) StepWin() uint64 {
	return f.stepWin
}
