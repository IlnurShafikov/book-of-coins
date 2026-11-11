// Package rules contain game rules, constants, and configuration structures
// that define how the game operates, including symbol IDs, paytables, and game mechanics.
package rules

// Symbol ID constants define unique identifiers for each symbol in the game.
// TODO: implement symbols ids according to your game design
const (
	// VoidSymbol represents an empty or null symbol position
	VoidSymbol = 0

	// High1 represents the first high-paying symbol
	High1 = 1

	// High2 represents the second high-paying symbol
	High2 = 2

	// Med1 represents the first medium-paying symbol
	Med1 = 3

	// Med2 represents the second medium-paying symbol
	Med2 = 4

	// Ace represents the Ace card symbol (low-paying)
	Ace = 5

	// King represents the King card symbol (low-paying)
	King = 6

	// Queen represents the Queen card symbol (low-paying)
	Queen = 7

	// Jack represents the Jack card symbol (low-paying)
	Jack = 8

	// Ten represents the Jack card symbol (low-paying)
	Ten = 9

	// ScatterID represents the scatter symbol that triggers bonus features
	ScatterID = 10

	// CoinID represents the regular coin symbol
	CoinID = 11
)

// Game mechanics constants
const (
	// InitRespinsCount defines the initial number of respins when the respin feature is triggered
	InitRespinsCount    = 4
	Step_TYPE_BASE_SPIN = 1
	Step_TYPE_FREE_SPIN = 2
	Step_TYPE_RE_SPIN   = 3
)
