package events

// ExpandSymbolEvent represents an expand symbol placement event.
// It tracks which reels received the special expanding symbol during a free spin.
type ExpandSymbolEvent struct {
	// Type identifies this event as a SpecSymbolEvent event.
	Type Type `json:"type"`
	// ExpandReels contains 0-based indices of reels where the symbol was placed.
	ExpandReels []int `json:"expandReels"`
	// Symbol is the ID of the expanding symbol that was placed.
	Symbol int `json:"symbol"`
}

// NewExpandSymbolEvent creates a new SpecSymbolEvent for the given expand symbol placement.
//
// Parameters:
//   - expandReels: 0-based indices of reels where the symbol was placed
//   - expandSymbolID: ID of the expanding symbol
//
// Returns:
//   - pointer to a new SpecSymbolEvent ready for client delivery
func NewExpandSymbolEvent(
	expandReels []int, expandSymbolID int,
) *ExpandSymbolEvent {
	return &ExpandSymbolEvent{
		Type:        ExpandSymbolType,
		ExpandReels: expandReels,
		Symbol:      expandSymbolID,
	}
}

// EventType returns the event type identifier, satisfying the [Event] interface.
//
// Returns:
//   - SpecSymbolEvent constant
func (s *ExpandSymbolEvent) EventType() Type {
	return s.Type
}
