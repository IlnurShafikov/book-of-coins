package result

import "github.com/IlnurShafikov/book-of-coins/internal/entities/result/events"

// ReSpin represents the result of a single respin.
// Respins are typically triggered by special symbols or mechanics during base or free spins.
type ReSpin struct {
	// boardResult is the board state before the respin was executed
	boardResult *BoardResult

	// events store all events that occurred during this respin
	events []events.Event

	// stepWin is the win amount for this specific respin
	stepWin uint64

	// totalWin is the cumulative win of the round up to and including this respin
	totalWin uint64

	// reTriggers is the number of additional respins triggered by this respin
	reTriggers uint32

	// maxWinReached indicates if the maximum win limit was reached
	maxWinReached bool
}

// NewReSpin creates a new ReSpin result with default values.
//
// Returns:
//   - pointer to a new ReSpin instance
func NewReSpin() *ReSpin {
	return &ReSpin{
		boardResult:   nil,
		events:        nil,
		totalWin:      0,
		stepWin:       0,
		reTriggers:    0,
		maxWinReached: false,
	}
}

// SetBoardResult sets the board result.
//
// Parameters:
//   - boardRes: pointer to the board result
func (r *ReSpin) SetBoardResult(boardRes *BoardResult) {
	r.boardResult = boardRes
}

// SetEvents sets the events that occurred during this respin.
//
// Parameters:
//   - e: slice of events to set
func (r *ReSpin) SetEvents(e []events.Event) {
	r.events = e
}

// SetReTriggers sets the number of additional respins triggered.
//
// Parameters:
//   - count: number of additional respins
func (r *ReSpin) SetReTriggers(count uint32) {
	r.reTriggers = count
}

// SetTotalWin sets the cumulative total win amount for this respin.
//
// Parameters:
//   - win: total win amount
func (r *ReSpin) SetTotalWin(win uint64) {
	r.totalWin = win
}

// SetStepWin sets the win amount for this specific respin.
//
// Parameters:
//   - win: win amount for this respin
func (r *ReSpin) SetStepWin(win uint64) {
	r.stepWin = win
}

// SetMaxWinReached sets whether the maximum win limit was reached.
//
// Parameters:
//   - reached: true if max win was reached, false otherwise
func (r *ReSpin) SetMaxWinReached(reached bool) {
	r.maxWinReached = reached
}

// BoardResult returns the board resul
//
// Returns:
//   - pointer to the board result before the respin
func (r *ReSpin) BoardResult() *BoardResult {
	return r.boardResult
}

// Events returns all events that occurred during this respin.
//
// Returns:
//   - slice of events
func (r *ReSpin) Events() []events.Event {
	return r.events
}

// ReTriggers returns the number of additional respins triggered.
//
// Returns:
//   - number of retriggers
func (r *ReSpin) ReTriggers() uint32 {
	return r.reTriggers
}

// StepWin returns the win amount for this specific respin.
//
// Returns:
//   - win amount for this respin
func (r *ReSpin) StepWin() uint64 {
	return r.stepWin
}

// TotalWin returns the cumulative total win amount for this respin.
//
// Returns:
//   - total win amount
func (r *ReSpin) TotalWin() uint64 {
	return r.totalWin
}

// MaxWinReached returns whether the maximum win limit was reached.
//
// Returns:
//   - true if max win was reached, false otherwise
func (r *ReSpin) MaxWinReached() bool {
	return r.maxWinReached
}
