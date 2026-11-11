// Package respin contains data transfer objects (DTOs) for respin phase.
// These structures represent respin results sent to the client.
package respin

// Step represents a complete respin step for client communication.
type Step struct {
	// Payload contains all the game data (symbols, events, counters)
	Payload Payload `json:"payload"`

	// TotalWin is the cumulative win amount up to and including this step
	TotalWin uint64 `json:"totalWin"`

	// StepWin is the win amount for this specific step
	StepWin uint64 `json:"stepWin"`

	// Type indicates the step type (should be TYPE_RE_SPIN)
	Type int `json:"type"`

	// MaxWinReached indicates if the maximum win limit was reached
	MaxWinReached bool `json:"maxWinReached"`
}
