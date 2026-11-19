package configs

import (
	"github.com/IlnurShafikov/book-of-coins/internal/entities/intervals"
	"github.com/IlnurShafikov/slot-engine/pkg/lines"
	"github.com/IlnurShafikov/slot-engine/pkg/paytable"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/fscounter"
)

// GeneralConfig holds the general game configuration that is shared across all game modes.
// It includes payout tables and various multipliers used throughout the game.
type GeneralConfig struct {
	// PayTable contains symbol payouts and payout ranges for the game
	PayTable *paytable.PayTable
	// ScatterPayTable contains the pay-anywhere payouts of the scatter symbol
	ScatterPayTable *paytable.PayTable
	// GrandRules describes the grand jackpot that is paid for a full board of coins
	GrandRules *GrandRules
	// JpTypes maps each jackpot value to its corresponding jackpot type identifier.
	JpTypes map[uint64]int
	// ExpandMinCount maps a special symbol ID to the minimum number of reels it must occupy
	ExpandMinCount map[int]int
	// FsInitCounter maps the number of scatter symbols to the number of free spins awarded.
	FsInitCounter fscounter.FSInitCounter
	// RsInitCounter maps the number of coins on the board to the number of respins
	RsInitCounter intervals.Intervals

	// Multipliers contains various multiplier configurations (cost, max win, etc.)
	Multipliers *MultipliersConfigs

	// SymbolsOnReels is the number of visible symbol rows per reel.
	SymbolsOnReels []int

	// Volatility represents the game's volatility setting, affecting risk and payout frequency.
	Volatility string

	// WinLines defines the cell-position patterns that form paying line combinations.
	WinLines lines.Lines
}

// GrandRules describes the grand jackpot of the respin phase.
type GrandRules struct {
	// CoinCountToGrand is the number of coins on a full board that pays the grand jackpot
	CoinCountToGrand int
	// GrandJpType is the jackpot type ID of the grand jackpot
	GrandJpType int
	// GrandWin is the grand jackpot prize in units of the cost multiplier
	GrandWin uint64
}

// MultipliersConfigs contains all multiplier-related configuration values.
type MultipliersConfigs struct {
	// MaxWinMultiplier defines the maximum win as a multiplier of the bet
	MaxWinMultiplier uint64

	// SimpleCostMultiplier defines the cost multiplier for the simple game mode
	SimpleCostMultiplier uint64
}

// MaxWinLimit calculates the absolute maximum win limit for the game.
// The limit is calculated as MaxWinMultiplier * SimpleCostMultiplier.
//
// Returns:
//   - the maximum win limit value
func (c *GeneralConfig) MaxWinLimit() uint64 {
	return c.Multipliers.MaxWinMultiplier * c.SimpleCostMultiplier()
}

// SimpleCostMultiplier returns the cost multiplier for the simple game mode.
//
// Returns:
//   - the simple mode cost multiplier
func (c *GeneralConfig) SimpleCostMultiplier() uint64 {
	return c.Multipliers.SimpleCostMultiplier
}
