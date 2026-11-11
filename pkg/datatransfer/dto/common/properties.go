package common

// SymbolProperties contains optional properties that can be attached to symbols.
// These properties modify the behavior or display of symbols.
type SymbolProperties struct {
	// Value is the coin value for symbols with monetary properties
	Value uint64 `json:"value"`

	// JpType indicates the jackpot type associated with the symbol (if any)
	JpType int `json:"jpType,omitempty"`
}
