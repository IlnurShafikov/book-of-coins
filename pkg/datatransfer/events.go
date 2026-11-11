package datatransfer

import "github.com/IlnurShafikov/book-of-coins/internal/entities/result/events"

// stepEventProvider is implemented by any game step that exposes its events.
type stepEventProvider interface {
	Events() []events.Event
}

// ExtractAndTransformEvents retrieves events from a step and scales monetary values by bet.
func ExtractAndTransformEvents(step stepEventProvider, betValue uint64) ([]events.Event, error) {
	eventList := step.Events()

	for _, curEvent := range eventList {
		if err := transformSingleEvent(curEvent, betValue); err != nil {
			return nil, err
		}
	}

	return eventList, nil
}

// transformSingleEvent routes a single event to its appropriate transformation handler.
//
// Parameters:
//   - curEvent: the event to transform
//   - betValue: the player's bet value
//
// Returns:
//   - error if the event type is unknown or transformation fails
func transformSingleEvent(curEvent events.Event, betValue uint64) error {
	switch curEvent.EventType() {
	case events.WinsType:
		return transformWinEvent(curEvent, betValue)
	case events.CollectCoinType:
		return transfomCollectCoinEvent(curEvent, betValue)
	case events.GrandWinType:
		return transfomGrandWinEvent(curEvent, betValue)
	case events.ExpandSymbolType:
		return nil
	default:
		return ErrInvalidEventType
	}
}

// transformWinEvent transforms a line wins event by multiplying payouts by bet value.
//
// Parameters:
//   - curEvent: the event to transform (must be *events.LineWinsEvent)
//   - betValue: the player's bet value
//
// Returns:
//   - error if the event is not the correct type
func transformWinEvent(curEvent events.Event, betValue uint64) error {
	winEvent, ok := curEvent.(*events.WinsEvent)
	if !ok {
		return ErrInvalidEventType
	}

	for i := range winEvent.LineWins {
		winEvent.LineWins[i].Payout *= betValue
	}

	for i := range winEvent.AnyWhereWins {
		winEvent.AnyWhereWins[i].Payout *= betValue
	}

	return nil
}

func transfomCollectCoinEvent(curEvent events.Event, betValue uint64) error {
	collectCoinEvent, ok := curEvent.(*events.CollectCoinEvent)
	if !ok {
		return ErrInvalidEventType
	}

	for i := range collectCoinEvent.Coins {
		collectCoinEvent.Coins[i].Value *= betValue
	}

	return nil
}

func transfomGrandWinEvent(curEvent events.Event, betValue uint64) error {
	grandWinEvent, ok := curEvent.(*events.GrandWinEvent)
	if !ok {
		return ErrInvalidEventType
	}

	grandWinEvent.Value *= betValue

	return nil
}
