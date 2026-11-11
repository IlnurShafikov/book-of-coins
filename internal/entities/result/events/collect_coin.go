package events

import "github.com/IlnurShafikov/book-of-coins/internal/entities/coins"

// CollectCoinEvent represents an event where coins are collected from the board.
// This event occurs when coin symbols are gathered for payout calculation.
type CollectCoinEvent struct {
	// Type is the event type identifier
	Type  Type             `json:"type"`
	Coins coins.CoinResult `json:"targets"`
}

// NewCollectCoinEvent creates a new coin collection event.
//
// Returns:
//   - pointer to a new CollectCoinEvent
func NewCollectCoinEvent(coinResult coins.CoinResult) *CollectCoinEvent {
	return &CollectCoinEvent{
		Type:  CollectCoinType,
		Coins: coinResult,
	}
}

// EventType returns the type identifier for this event.
func (c *CollectCoinEvent) EventType() Type {
	return c.Type
}
