package result

// RespinGame accumulates statistics for the respin phase, split by the phase that triggered it.
type RespinGame struct {
	// WinsFromBase is the sum of the wins of respins triggered by the base game.
	WinsFromBase uint64
	// WinsFromFree is the sum of the wins of respins triggered by free spins.
	WinsFromFree uint64
	// ReTriggersCountFromBase is the number of base game respins that restored the counter.
	ReTriggersCountFromBase int
	// ReTriggersCountFromFree is the number of free game respins that restored the counter.
	ReTriggersCountFromFree int
	// ReTriggerStepsCountFromBase is the number of respins added by base game retriggers.
	ReTriggerStepsCountFromBase uint32
	// ReTriggerStepsCountFromFree is the number of respins added by free game retriggers.
	ReTriggerStepsCountFromFree uint32
}

// newRespinGame creates a RespinGame with zeroed statistics.
func newRespinGame() *RespinGame {
	return &RespinGame{
		WinsFromBase:                0,
		WinsFromFree:                0,
		ReTriggerStepsCountFromBase: 0,
		ReTriggerStepsCountFromFree: 0,
		ReTriggersCountFromBase:     0,
		ReTriggersCountFromFree:     0,
	}
}

// AddWinsFromBase adds the win of a respin of the base game.
func (r *RespinGame) AddWinsFromBase(wins uint64) {
	r.WinsFromBase += wins
}

// AddWinsFromFree adds the win of a respin of a free spin.
func (r *RespinGame) AddWinsFromFree(wins uint64) {
	r.WinsFromFree += wins
}

// AddReTriggersCountFromBase records a respin of the base game that restored the counter.
// Respins without retrigger (0) are ignored.
//
// Parameters:
//   - reTriggersCount: number of respins added by the retrigger
func (r *RespinGame) AddReTriggersCountFromBase(reTriggersCount uint32) {
	if reTriggersCount > 0 {
		r.ReTriggerStepsCountFromBase += reTriggersCount
		r.ReTriggersCountFromBase++
	}
}

// AddReTriggersCountFromFree records a respin of a free spin that restored the counter.
// Respins without retrigger (0) are ignored.
//
// Parameters:
//   - reTriggersCount: number of respins added by the retrigger
func (r *RespinGame) AddReTriggersCountFromFree(reTriggersCount uint32) {
	if reTriggersCount > 0 {
		r.ReTriggerStepsCountFromFree += reTriggersCount
		r.ReTriggersCountFromFree++
	}
}

// GetWinsFromBase returns the sum of the wins of respins triggered by the base game.
func (r *RespinGame) GetWinsFromBase() uint64 {
	return r.WinsFromBase
}

// GetWinsFromFree returns the sum of the wins of respins triggered by free spins.
func (r *RespinGame) GetWinsFromFree() uint64 {
	return r.WinsFromFree
}

// GetReTriggerStepsCountFromBase returns the number of respins added by base game retriggers.
func (r *RespinGame) GetReTriggerStepsCountFromBase() uint32 {
	return r.ReTriggerStepsCountFromBase
}

// GetReTriggerStepsCountFromFree returns the number of respins added by free game retriggers.
func (r *RespinGame) GetReTriggerStepsCountFromFree() uint32 {
	return r.ReTriggerStepsCountFromFree
}

// GetReTriggersCountFromBase returns the number of base game respins that restored the counter.
func (r *RespinGame) GetReTriggersCountFromBase() int {
	return r.ReTriggersCountFromBase
}

// GetReTriggersCountFromFree returns the number of free game respins that restored the counter.
func (r *RespinGame) GetReTriggersCountFromFree() int {
	return r.ReTriggersCountFromFree
}

// Merge adds another RespinGame's statistics into this one.
func (r *RespinGame) Merge(other *RespinGame) {
	r.WinsFromBase += other.WinsFromBase
	r.WinsFromFree += other.WinsFromFree
	r.ReTriggersCountFromBase += other.ReTriggersCountFromBase
	r.ReTriggersCountFromFree += other.ReTriggersCountFromFree
	r.ReTriggerStepsCountFromBase += other.ReTriggerStepsCountFromBase
	r.ReTriggerStepsCountFromFree += other.ReTriggerStepsCountFromFree
}
