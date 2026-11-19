// Package parser provides the RoundParser which processes raw game round results
// and accumulates them into a simulation Result for subsequent statistical analysis.
package parser

import (
	resp "github.com/IlnurShafikov/book-of-coins/internal/entities/result"
	"github.com/IlnurShafikov/book-of-coins/internal/simulator/result"
)

// RoundParser processes round results and accumulates simulation statistics
type RoundParser struct {
	// res accumulates the statistics of all processed rounds.
	res *result.Result
}

// NewRoundParser creates a new RoundParser with initialized result container
func NewRoundParser(roundCost uint64, simpleCost int) *RoundParser {
	return &RoundParser{
		res: result.New(roundCost, simpleCost),
	}
}

// MergeResult merges a worker's accumulated result into this parser's result.
// Used by the parallel runner after all workers complete.
func (rp *RoundParser) MergeResult(other *result.Result) {
	rp.res.Merge(other)
}

// Parse processes a complete round result and updates simulation statistics
func (rp *RoundParser) Parse(roundResult *resp.Round) error {
	var (
		totalRoundWin               uint64
		totalFreeWin                uint64
		totalRespinFromFreeTotalWin uint64
	)

	totalBaseWin, totalRespinFromBaseWin := rp.parseBasePhase(roundResult.BasePhase)

	if len(roundResult.FreePhase) > 0 {
		totalFreeWin, totalRespinFromFreeTotalWin = rp.parseFreePhase(roundResult.FreePhase)

		rp.res.Base().AddFreeGameSteps(len(roundResult.FreePhase))
	}

	totalRoundWin = totalBaseWin + totalFreeWin +
		totalRespinFromBaseWin + totalRespinFromFreeTotalWin

	rp.updateRoundStatistics(totalRoundWin)

	return nil
}

// parseBasePhase processes base game phase results
func (rp *RoundParser) parseBasePhase(
	baseResult *resp.Base,
) (baseTotalWin, totalRespinWin uint64) {
	baseTotalWin = baseResult.TotalWin()

	rp.res.Base().AddWins(baseTotalWin)
	rp.res.BaseWinRanges.AddWinRange(baseTotalWin)

	if baseTotalWin > 0 {
		rp.res.BasePositiveBets().Add()
	}

	if baseResult.MaxWinReached() {
		rp.res.AddMaxWinReached()
	}

	if len(baseResult.ReSpins) > 0 {
		totalRespinWin = rp.parseRespinPhaseFromBase(baseResult.ReSpins)
		rp.res.BaseReSpinPositiveBets().Add()
		rp.res.BaseReSpinWinRanges.AddWinRange(totalRespinWin)
		rp.res.Base().AddRespinGameSteps(len(baseResult.ReSpins))
	}

	return baseTotalWin, totalRespinWin
}

// parseRespinPhaseFromFree processes the respins triggered by a free spin.
func (rp *RoundParser) parseRespinPhaseFromFree(respinSteps []*resp.ReSpin) uint64 {
	return rp.parseRespinPhase(
		respinSteps,
		rp.processSingleRespinStepFromFree,
	)
}

// parseFreePhase processes all free spin results
func (rp *RoundParser) parseFreePhase(
	freeSteps []*resp.Free,
) (totalFreeWin, respinTotalWin uint64) {
	for _, freeStep := range freeSteps {
		fsStepWin, respinTotalWinInFs := rp.processSingleFreeStep(freeStep)
		totalFreeWin += fsStepWin
		respinTotalWin += respinTotalWinInFs
	}

	rp.res.FreeWinRanges.AddWinRange(totalFreeWin)
	rp.res.FreeReSpinWinRanges.AddWinRange(respinTotalWin)

	if totalFreeWin > 0 {
		rp.res.FreePositiveBets().Add()
	}

	if respinTotalWin > 0 {
		rp.res.FreeReSpinPositiveBets().Add()
	}

	return totalFreeWin, respinTotalWin
}

// processSingleFreeStep processes a single free spin step
func (rp *RoundParser) processSingleFreeStep(
	freeResult *resp.Free,
) (stepWin, respinTotalWin uint64) {
	stepWin = freeResult.StepWin()

	rp.res.Free().AddWins(stepWin)

	if freeResult.MaxWinReached() {
		rp.res.AddMaxWinReached()
	}

	if len(freeResult.ReSpins) > 0 {
		respinTotalWin = rp.parseRespinPhaseFromFree(freeResult.ReSpins)
		rp.res.Free().AddRespinGameSteps(len(freeResult.ReSpins))
	}

	return stepWin, respinTotalWin
}

// parseRespinPhaseFromBase processes the respins triggered by the base spin.
func (rp *RoundParser) parseRespinPhaseFromBase(respinSteps []*resp.ReSpin) uint64 {
	return rp.parseRespinPhase(
		respinSteps,
		rp.processSingleRespinStepFromBase,
	)
}

// parseRespinPhase sums the wins of the respins using the given step processor.
//
// Parameters:
//   - respinSteps: respins to process
//   - processStep: collects the statistics of one respin and returns its win
//
// Returns:
//   - total win of all respins
func (rp *RoundParser) parseRespinPhase(
	respinSteps []*resp.ReSpin,
	processStep func(*resp.ReSpin) uint64,
) uint64 {
	var totalRespinWin uint64

	for _, respinStep := range respinSteps {
		stepWin := processStep(respinStep)
		totalRespinWin += stepWin
	}

	return totalRespinWin
}

// processSingleRespinStepFromBase collects the statistics of a base game respin
// and returns its win.
func (rp *RoundParser) processSingleRespinStepFromBase(respinResult *resp.ReSpin) uint64 {
	stepWin := respinResult.StepWin()

	rp.res.ReSpin().AddWinsFromBase(stepWin)
	rp.res.ReSpin().AddReTriggersCountFromBase(respinResult.ReTriggers())

	if respinResult.MaxWinReached() {
		rp.res.AddMaxWinReached()
	}

	return stepWin
}

// processSingleRespinStepFromFree collects the statistics of a free game respin
// and returns its win.
func (rp *RoundParser) processSingleRespinStepFromFree(respinResult *resp.ReSpin) uint64 {
	stepWin := respinResult.StepWin()

	rp.res.ReSpin().AddWinsFromFree(stepWin)
	rp.res.ReSpin().AddReTriggersCountFromFree(respinResult.ReTriggers())

	if respinResult.MaxWinReached() {
		rp.res.AddMaxWinReached()
	}

	return stepWin
}

// updateRoundStatistics updates overall round statistics
func (rp *RoundParser) updateRoundStatistics(totalWin uint64) {
	if totalWin > 0 {
		rp.res.Xi().Add(totalWin)
		rp.res.AddMaxTotalWin(totalWin)
		rp.res.PositiveBets().Add()
	}
}

// GetResult returns the accumulated simulation results
func (rp *RoundParser) GetResult() *result.Result {
	return rp.res
}
