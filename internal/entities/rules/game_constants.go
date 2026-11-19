// Package rules contain game rules, constants, and configuration structures
// that define how the game operates, including symbol IDs, paytables, and game mechanics.
package rules

// Symbol ID constants define unique identifiers for each symbol in the game.
// The IDs match the reel strips, the paytable and the special symbol weights.
const (
	// VoidSymbol represents an empty position; it is used for the empty cells of the
	// respin phase and never appears on the reel strips
	VoidSymbol = 0

	// High1 is the most valuable regular symbol
	High1 = 1

	// High2 is the second high-pay symbol
	High2 = 2

	// High3 is the third high-pay symbol
	High3 = 3

	// High4 is the fourth high-pay symbol
	High4 = 4

	// Low1 is the first low-pay symbol
	Low1 = 5

	// Low2 is the second low-pay symbol
	Low2 = 6

	// Low3 is the third low-pay symbol
	Low3 = 7

	// Low4 is the fourth low-pay symbol
	Low4 = 8

	// Low5 is the least valuable symbol
	Low5 = 9

	// ScatterID is the Wild/Scatter symbol: one symbol that substitutes for the regular
	// symbols on the lines and triggers the free spins when it lands anywhere on the board
	ScatterID = 10

	// CoinID is the coin symbol that carries a cash or jackpot value and triggers the respins
	CoinID = 11
)

// Game mechanics constants
const (
	// InitRespinsCount is the number of respins at the start of the respin phase and the
	// value the counter is restored to every time a new coin lands
	InitRespinsCount = 4

	// Step_TYPE_BASE_SPIN is the type of a base spin step.
	Step_TYPE_BASE_SPIN = 1
	// Step_TYPE_FREE_SPIN is the type of a free spin step.
	Step_TYPE_FREE_SPIN = 2
	// Step_TYPE_RE_SPIN is the type of a respin step.
	Step_TYPE_RE_SPIN = 3
)
