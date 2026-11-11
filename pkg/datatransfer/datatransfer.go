package datatransfer

import (
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/entities/result"
	"github.com/IlnurShafikov/book-of-coins/pkg/datatransfer/dto/common"
	"github.com/IlnurShafikov/book-of-coins/pkg/datatransfer/dto/step"
)

// Transfer is responsible for transforming internal game results into client-facing DTOs.
// It handles the conversion of all game phases (base, free, respins) into API responses.
type Transfer struct {
	// costMultiplier is the base cost multiplier for the game mode
	costMultiplier uint64
}

// NewTransfer creates a new Transfer instance with the specified cost multiplier.
//
// Parameters:
//   - costMultiplier: the base cost multiplier for game rounds
//
// Returns:
//   - pointer to a new Transfer instance
func NewTransfer(
	costMultiplier uint64,
) *Transfer {
	resp := &Transfer{
		costMultiplier: costMultiplier,
	}

	return resp
}

// Transfer converts a complete round result into a BetResult for the API.
// It orchestrates the transformation of base spin, free spins, and respins into client steps.
//
// Parameters:
//   - res: the complete round result from the game engine
//   - betValue: the player's bet amount
//
// Returns:
//   - pointer to BetResult containing all transformed steps
//   - error if transformation fails
func (t *Transfer) Transfer(res *result.Round, betValue uint64) (*step.BetResult, error) {
	var (
		totalRoundWin uint64
		err           error
	)

	// Pre-allocate slice with estimated capacity for all steps
	allCliSteps := make([]*step.Step, 0, countStepsCapacity(res)+1)
	roundCost := res.CostMultiplier * betValue

	// Transform the base game spin
	allCliSteps, totalRoundWin, err = t.TransformBaseStep(
		res.BasePhase, allCliSteps, roundCost, betValue)
	if err != nil {
		return nil, fmt.Errorf("transfer base step: %w", err)
	}

	// baseParent used when we get maxWin in other steps to return in correct state
	baseParent, err := defaultParentStep(res.BasePhase.BoardResult(), betValue)
	if err != nil {
		return nil, fmt.Errorf("create base parent step: %w", err)
	}

	// Transform free spins if any were triggered
	allCliSteps, totalRoundWin, err = t.TransformFreeStep(
		res, allCliSteps, baseParent, totalRoundWin, betValue)
	if err != nil {
		return nil, fmt.Errorf("transfer free steps: %w", err)
	}

	// Transform respins from base phase if any were triggered
	allCliSteps, _, err = t.TransformReSpinSteps(
		res.BasePhase.ReSpins, allCliSteps, baseParent, baseParent,
		res.BasePhase.ReSpinInitials(), totalRoundWin, betValue)
	if err != nil {
		return nil, fmt.Errorf("transfer respins: %w", err)
	}

	// Construct the final bet result
	clientRoundResp := &step.BetResult{
		Steps:          allCliSteps,
		UserBet:        betValue,
		CostMultiplier: t.costMultiplier,
		Seed:           res.Seed,
	}

	return clientRoundResp, nil
}

// countStepsCapacity estimates the total number of steps in the round for efficient allocation.
// This helps optimize memory allocation by pre-sizing the steps slice.
//
// Parameters:
//   - res: the round result to count steps from
//
// Returns:
//   - estimated number of steps in the round
func countStepsCapacity(res *result.Round) int {
	var totalSteps int

	totalSteps += len(res.BasePhase.ReSpins)
	totalSteps += len(res.FreePhase)

	for _, freeStep := range res.FreePhase {
		totalSteps += len(freeStep.ReSpins)
	}

	return totalSteps
}

// defaultParentStep creates a default ParentStep with transformed symbols
// for use in max win scenarios.
func defaultParentStep(gameBoard *result.BoardResult, betValue uint64) (*common.ParentStep, error) {
	parentSymbols, err := TransformSymbols(gameBoard, betValue)
	if err != nil {
		return nil, fmt.Errorf("transform symbols: %w", err)
	}

	return &common.ParentStep{
		Payload: common.ParentPayload{
			Counter: nil,
			Symbols: parentSymbols,
		},
	}, nil
}
