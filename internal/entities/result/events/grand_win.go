package events

// GrandWinEvent is sent when the grand jackpot is paid for a full board of coins.
type GrandWinEvent struct {
	// Type is the event type identifier
	Type Type `json:"type"`

	// Value is the monetary value of the coin
	Value uint64 `json:"value"`

	// JpType is the jackpot type identifier (0 if not a jackpot coin)
	JpType int `json:"jpType"`
}

// NewGrandWinEvent creates a grand win event.
//
// Parameters:
//   - value: grand jackpot prize
//   - jpType: jackpot type ID of the grand jackpot
//
// Returns:
//   - pointer to a new GrandWinEvent
func NewGrandWinEvent(value uint64, jpType int) *GrandWinEvent {
	return &GrandWinEvent{
		Type:   GrandWinType,
		Value:  value,
		JpType: jpType,
	}
}

// EventType returns the type identifier for this event.
func (c *GrandWinEvent) EventType() Type {
	return c.Type
}
