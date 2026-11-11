package respin

import (
	"github.com/IlnurShafikov/book-of-coins/internal/entities/result/events"
	"github.com/IlnurShafikov/book-of-coins/pkg/datatransfer/dto/common"
)

// Payload contains the detailed game data for a respin step.
type Payload struct {
	// ParentStep contains the parent base data
	ParentStep *common.ParentStep `json:"parent,omitempty"`

	// Symbols represents the current symbol board
	Symbols []*common.Symbol `json:"symbols"`

	// Events contains all game events that occurred (wins, retriggers, etc.)
	Events []events.Event `json:"events,omitempty"`

	// Counter tracks the respins (added, left, played)
	Counter common.Counter `json:"counter"`

	// NextStepType indicates what type of step comes next (respin, free, or base)
	NextStepType int `json:"nextStepType"`
}
