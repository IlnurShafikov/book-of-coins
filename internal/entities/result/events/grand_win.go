package events

type GrandWinEvent struct {
	// Type is the event type identifier
	Type Type `json:"type"`

	// Value is the monetary value of the coin
	Value uint64 `json:"value"`

	// JpType is the jackpot type identifier (0 if not a jackpot coin)
	JpType int `json:"jpType"`
}

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
