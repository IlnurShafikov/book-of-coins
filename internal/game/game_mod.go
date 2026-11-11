package game

import (
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/components/limiter"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/result"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
)

// freeSpinsCapacityFactor pre-sizes the free spins slice to hold the initial
// spins plus a typical number of retriggers.
const freeSpinsCapacityFactor = 2

// BasePhase defines the interface for base game phase operations
type BasePhase interface {
	Spin(rng model.RNG, winLimiter limiter.WinLimiter) (*result.Base, error)
}

// FreePhase defines the interface for free spins phase operations
type FreePhase interface {
	Spin(
		rng model.RNG,
		expandSymbolID int,
		winLimiter limiter.WinLimiter,
	) (*result.Free, error)
}

type ReSpinPhase interface {
	Spins(
		rng model.RNG,
		rsCount int,
		prevBoardResult *result.BoardResult,
		winLimiter limiter.WinLimiter,
	) ([]*result.ReSpin, error)
}

// Mod represents a complete game modification with base and free spin phases.
// It coordinates all game phases and manages win limits and cost multipliers.
type Mod struct {
	// baseSpinPhase handles base game spin logic
	baseSpinPhase BasePhase

	// freeSpinPhase handles free spins logic
	freeSpinPhase FreePhase

	// reSpinPhase handles respin logic
	reSpinPhase ReSpinPhase

	// costMultiplier is the multiplier applied to bet values
	costMultiplier uint64

	// maxWinLimit is the maximum win amount allowed per round
	maxWinLimit uint64
}

// NewMod creates a new Mod instance with the specified phases and limits.
//
// Parameters:
//   - basePhase: base game phase implementation
//   - freeSpinPhase: free spins phase implementation
//   - reSpinPhase: respin phase implementation
//   - costMultiplier: multiplier for bet values
//   - maxWinLimit: maximum win amount allowed per round
//
// Returns:
//   - pointer to the newly created Mod
func NewMod(
	basePhase BasePhase,
	freeSpinPhase FreePhase,
	reSpinPhase ReSpinPhase,
	costMultiplier, maxWinLimit uint64,
) *Mod {
	return &Mod{
		baseSpinPhase:  basePhase,
		freeSpinPhase:  freeSpinPhase,
		reSpinPhase:    reSpinPhase,
		costMultiplier: costMultiplier,
		maxWinLimit:    maxWinLimit,
	}
}

// CostMultiplier returns the cost multiplier for this game mod
func (g *Mod) CostMultiplier() uint64 {
	return g.costMultiplier
}

// Spin executes a complete game round and returns the result
func (g *Mod) Spin(rng model.RNG) (*result.Round, error) {
	winLimiter, err := limiter.New(g.maxWinLimit)
	if err != nil {
		return nil, fmt.Errorf("initialize maxWinCollector %w", err)
	}

	roundResult := result.NewRound(g.costMultiplier)

	if err := g.startBasePhase(rng, roundResult, winLimiter); err != nil {
		return nil, fmt.Errorf("execute base phase: %w", err)
	}

	if err := g.startFreePhaseIfTriggered(rng, roundResult, winLimiter); err != nil {
		return nil, fmt.Errorf("execute free phase: %w", err)
	}

	return roundResult, nil
}

// startBasePhase executes the base game phase and updates the round result
func (g *Mod) startBasePhase(
	rng model.RNG, roundResult *result.Round, winLimiter limiter.WinLimiter,
) error {
	baseResult, err := g.baseSpinPhase.Spin(rng, winLimiter)
	if err != nil {
		return fmt.Errorf("spin base game: %w", err)
	}

	err = g.baseStartReSpinPhaseIfTriggered(rng, baseResult, winLimiter)
	if err != nil {
		return fmt.Errorf("start re spin phase: %w", err)
	}

	roundResult.BasePhase = baseResult

	return nil
}

func (g *Mod) startFreePhaseIfTriggered(
	rng model.RNG, roundResult *result.Round, winLimiter limiter.WinLimiter,
) error {
	fsCount := roundResult.BasePhase.FreeSpinsInitialsCount()
	if fsCount == 0 {
		return nil // No free spins triggered
	}

	allRes := make([]*result.Free, 0, fsCount*freeSpinsCapacityFactor)

	expandSymbolID := roundResult.BasePhase.ExpandSymbolID()

	// fsCount grows on retriggers, so the bound must be re-evaluated on every iteration.
	//nolint:intrange // range over int evaluates the bound only once
	for fsIndex := 0; fsIndex < int(fsCount); fsIndex++ {
		if winLimiter.MaxWinReached() {
			break
		}

		freeResult, err := g.freeSpinPhase.Spin(rng, expandSymbolID, winLimiter)
		if err != nil {
			return fmt.Errorf("spin free games (count: %d): %w", fsIndex, err)
		}

		fsCount += freeResult.ReTriggers()

		err = g.freeStartReSpinPhaseIfTriggered(rng, freeResult, winLimiter)
		if err != nil {
			return fmt.Errorf("start re spin phase: %w", err)
		}

		allRes = append(allRes, freeResult)
	}

	roundResult.FreePhase = allRes

	return nil
}

// baseStartReSpinPhaseIfTriggered checks for respin triggers from base phase and executes respins.
//
// Parameters:
//   - rng: random number generator
//   - baseResult: base phase result that may have triggered respins
//   - winLimiter: win limiter to cap maximum wins
//
// Returns:
//   - error if respin execution fails
func (g *Mod) baseStartReSpinPhaseIfTriggered(
	rng model.RNG,
	baseResult *result.Base,
	winLimiter limiter.WinLimiter,
) error {
	if baseResult.ReSpinInitials() == 0 || winLimiter.MaxWinReached() {
		return nil
	}

	reSpinRes, err := g.reSpinPhase.Spins(
		rng, int(baseResult.ReSpinInitials()), baseResult.BoardResult().Copy(), winLimiter)
	if err != nil {
		return fmt.Errorf("spin respin phase: %w", err)
	}

	if len(reSpinRes) > 0 {
		baseResult.ReSpins = reSpinRes
	}

	return nil
}

// freeStartReSpinPhaseIfTriggered checks for respin triggers from free spins and executes respins.
//
// Parameters:
//   - rng: random number generator
//   - freeResult: free spin result that may have triggered respins
//   - winLimiter: win limiter to cap maximum wins
//
// Returns:
//   - error if respin execution fails
func (g *Mod) freeStartReSpinPhaseIfTriggered(
	rng model.RNG,
	freeResult *result.Free,
	winLimiter limiter.WinLimiter,
) error {
	if freeResult.ReSpinInitials() == 0 || winLimiter.MaxWinReached() {
		return nil
	}

	reSpinRes, err := g.reSpinPhase.Spins(
		rng, int(freeResult.ReSpinInitials()), freeResult.AfterBoard().Copy(), winLimiter)
	if err != nil {
		return fmt.Errorf("spin respin phase: %w", err)
	}

	if len(reSpinRes) > 0 {
		freeResult.SetTotalWin(winLimiter.TotalWin())
		freeResult.ReSpins = reSpinRes
	}

	return nil
}
