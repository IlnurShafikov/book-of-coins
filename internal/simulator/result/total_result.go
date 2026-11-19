// Package result contains data structures that accumulate simulation statistics
// across all simulated rounds, including win totals, hit rates, RTP metrics,
// and win-range distributions for both base and free game phases.
package result

// Result holds the aggregated statistics collected during a simulation run.
// It is populated by the parser package as each round result is processed.
type Result struct {
	// BaseWinRanges tracks the distribution of base-phase total wins.
	BaseWinRanges *WinRanges
	// FreeWinRanges tracks the distribution of free-phase total wins.
	FreeWinRanges *WinRanges
	// BaseReSpinWinRanges tracks the distribution of respin wins after the base game.
	BaseReSpinWinRanges *WinRanges
	// FreeReSpinWinRanges tracks the distribution of respin wins after free spins.
	FreeReSpinWinRanges *WinRanges

	// base accumulates base game statistics.
	base *BaseGame
	// free accumulates free game statistics.
	free *FreeGame
	// respin accumulates respin statistics.
	respin *RespinGame

	// positiveBets counts rounds with any win.
	positiveBets *PositiveBets
	// baseGamePositiveBets counts rounds with a win in the base game.
	baseGamePositiveBets *PositiveBets
	// freeGamePositiveBets counts rounds with a win in the free game.
	freeGamePositiveBets *PositiveBets
	// freeReSpinGamePositiveBets counts rounds with a respin win after free spins.
	freeReSpinGamePositiveBets *PositiveBets
	// baseReSpinGamePositiveBets counts rounds with a respin win after the base game.
	baseReSpinGamePositiveBets *PositiveBets
	// xi accumulates squared RTP values for the volatility calculation.
	xi *Xi
	// maxTotalWin is the largest total win of a single round.
	maxTotalWin uint64
	// maxWinReached is the number of rounds that hit the max win limit.
	maxWinReached uint64
}

// New creates a new Result with all counters zeroed and win-range buckets initialised.
func New(roundCost uint64, simpleCostMultiplier int) *Result {
	return &Result{
		base:                       newBaseGame(),
		free:                       newFreeGame(),
		respin:                     newRespinGame(),
		positiveBets:               newPositiveBets(),
		baseGamePositiveBets:       newPositiveBets(),
		freeGamePositiveBets:       newPositiveBets(),
		freeReSpinGamePositiveBets: newPositiveBets(),
		baseReSpinGamePositiveBets: newPositiveBets(),
		xi:                         newXi(roundCost),
		maxTotalWin:                0,
		maxWinReached:              0,
		BaseWinRanges:              NewWinRanges(simpleCostMultiplier),
		FreeWinRanges:              NewWinRanges(simpleCostMultiplier),
		BaseReSpinWinRanges:        NewWinRanges(simpleCostMultiplier),
		FreeReSpinWinRanges:        NewWinRanges(simpleCostMultiplier),
	}
}

// Merge incorporates another Result into this one.
// Used by the parallel runner to combine per-worker accumulators.
func (r *Result) Merge(other *Result) {
	r.base.Merge(other.base)
	r.free.Merge(other.free)
	r.respin.Merge(other.respin)
	r.positiveBets.AddN(other.positiveBets.GetSpins())
	r.xi.AddSum(other.xi.Sum())
	r.maxWinReached += other.maxWinReached

	if other.maxTotalWin > r.maxTotalWin {
		r.maxTotalWin = other.maxTotalWin
	}

	r.BaseWinRanges.Merge(other.BaseWinRanges)
	r.FreeWinRanges.Merge(other.FreeWinRanges)
	r.BaseReSpinWinRanges.Merge(other.BaseReSpinWinRanges)
	r.FreeReSpinWinRanges.Merge(other.FreeReSpinWinRanges)
}

// Base returns the base game statistics collector.
func (r *Result) Base() *BaseGame { return r.base }

// Free returns the free game statistics collector.
func (r *Result) Free() *FreeGame { return r.free }

// ReSpin returns the respin statistics collector.
func (r *Result) ReSpin() *RespinGame { return r.respin }

// Xi returns the Xi statistics tracker.
func (r *Result) Xi() *Xi { return r.xi }

// MaxTotalWin returns the maximum total win across all rounds.
func (r *Result) MaxTotalWin() uint64 { return r.maxTotalWin }

// PositiveBets returns the positive bets statistics collector.
func (r *Result) PositiveBets() *PositiveBets { return r.positiveBets }

// BasePositiveBets returns the collector of base game rounds with a win.
func (r *Result) BasePositiveBets() *PositiveBets { return r.baseGamePositiveBets }

// FreePositiveBets returns the collector of free game rounds with a win.
func (r *Result) FreePositiveBets() *PositiveBets { return r.freeGamePositiveBets }

// FreeReSpinPositiveBets returns the collector of rounds with a respin win after free spins.
func (r *Result) FreeReSpinPositiveBets() *PositiveBets {
	return r.freeReSpinGamePositiveBets
}

// BaseReSpinPositiveBets returns the collector of rounds with a respin win after the base game.
func (r *Result) BaseReSpinPositiveBets() *PositiveBets {
	return r.baseReSpinGamePositiveBets
}

// AddMaxTotalWin updates the maximum total win if the provided win amount is greater.
func (r *Result) AddMaxTotalWin(win uint64) {
	if r.maxTotalWin < win {
		r.maxTotalWin = win
	}
}

// AddMaxWinReached increments the counter for maximum wins reached.
func (r *Result) AddMaxWinReached() {
	r.maxWinReached++
}

// MaxWinReached returns the total count of maximum wins reached.
func (r *Result) MaxWinReached() uint64 {
	return r.maxWinReached
}
