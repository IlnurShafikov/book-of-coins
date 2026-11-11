// Package rules provides the domain types and constructors for game rules
// that are serialised to JSON and sent to the client on game load.
// Rules include the paytable, volatility level, RTP value, and maximum win multiplier.
package rules

import (
	"encoding/json"
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/configs"
	"github.com/IlnurShafikov/slot-engine/pkg/lines"
	jsoniter "github.com/json-iterator/go"
)

// Rules represents the complete game rules sent to the client.
// It includes all information needed for the client to understand the game mechanics and payouts.
// TODO: implement the required fields for your specific game
type Rules struct {
	// Volatility describes the game's volatility level (risk/variance)
	Volatility *Volatility `json:"volatility"`

	// PayTable contains symbol payouts and match ranges
	PayTable *PayTable `json:"payTable"`

	// ScatterPayTable contains payouts for scatter (pay anywhere) symbols
	ScatterPayTable *PayTable `json:"scatterPaytable"`

	Lines lines.Lines `json:"lines"`

	JpTypes JpTypes `json:"jpTypes"`

	// MaxMultiplier is the maximum win multiplier relative to bet
	MaxMultiplier uint64 `json:"maxMultiplier"`

	// RTP is the expected Return to Player percentage
	RTP float64 `json:"rtp"`
}

// JpTypes maps each jackpot type ID to its payout multiplier.
// Serialised as a JSON object with integer string keys.
type JpTypes map[int]uint64

// MakeGameRules constructs the game rules JSON from configuration data.
// This JSON is sent to the client to inform them of the game's mechanics and payouts.
//
// Parameters:
//   - general: general configuration containing paytable and multipliers
//   - expRTP: expected RTP percentage for this game version
//
// Returns:
//   - json.RawMessage containing the serialized game rules
//   - error if marshaling fails
func MakeGameRules(
	general *configs.GeneralConfig,
	expRTP float64,
) (json.RawMessage, error) {
	paytable := NewPaytable(general.PayTable)
	scatterPaytable := NewPaytable(general.ScatterPayTable)
	jpTypes := makeJackpotMultiplier(general.JpTypes)

	rules := Rules{
		PayTable:        paytable,
		ScatterPayTable: scatterPaytable,
		Lines:           general.WinLines,
		JpTypes:         jpTypes,
		MaxMultiplier:   general.Multipliers.MaxWinMultiplier,
		RTP:             expRTP,
		Volatility:      &Volatility{Level: general.Volatility},
	}

	rulesJSON, err := jsoniter.Marshal(rules)
	if err != nil {
		return nil, fmt.Errorf("marshal rules: %w", err)
	}

	return rulesJSON, nil
}

// Volatility represents the game's volatility (variance/risk) level.
// Higher volatility means larger but less frequent wins,
// lower volatility means smaller but more frequent wins.
type Volatility struct {
	// Level is the volatility classification (e.g., "low", "medium", "high")
	Level string `json:"level"`
}

// makeJackpotMultiplier converts a Pairs slice of jackpot type → multiplier
// entries into the JpTypes map expected by the rules JSON payload.
//
// Parameters:
//   - jackpotPair: key-value pairs mapping jackpot type IDs to their multipliers
//
// Returns:
//   - JpTypes map with each jackpot type ID mapped to its payout multiplier
func makeJackpotMultiplier(
	jackpotTypes map[uint64]int,
) JpTypes {
	jpTypes := make(JpTypes, len(jackpotTypes))

	for value, jpType := range jackpotTypes {
		jpTypes[jpType] = value
	}

	return jpTypes
}
