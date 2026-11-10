package result

type RespinGame struct {
	WinsFromBase                uint64
	WinsFromFree                uint64
	ReTriggersCountFromBase     int
	ReTriggersCountFromFree     int
	ReTriggerStepsCountFromBase uint32
	ReTriggerStepsCountFromFree uint32
}

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

func (r *RespinGame) AddWinsFromBase(wins uint64) {
	r.WinsFromBase += wins
}

func (r *RespinGame) AddWinsFromFree(wins uint64) {
	r.WinsFromFree += wins
}

func (r *RespinGame) AddReTriggersCountFromBase(reTriggersCount uint32) {
	if reTriggersCount > 0 {
		r.ReTriggerStepsCountFromBase += reTriggersCount
		r.ReTriggersCountFromBase++
	}
}

func (r *RespinGame) AddReTriggersCountFromFree(reTriggersCount uint32) {
	if reTriggersCount > 0 {
		r.ReTriggerStepsCountFromFree += reTriggersCount
		r.ReTriggersCountFromFree++
	}
}

func (r *RespinGame) GetWinsFromBase() uint64 {
	return r.WinsFromBase
}

func (r *RespinGame) GetWinsFromFree() uint64 {
	return r.WinsFromFree
}

func (r *RespinGame) GetReTriggerStepsCountFromBase() uint32 {
	return r.ReTriggerStepsCountFromBase
}

func (r *RespinGame) GetReTriggerStepsCountFromFree() uint32 {
	return r.ReTriggerStepsCountFromFree
}

func (r *RespinGame) GetReTriggersCountFromBase() int {
	return r.ReTriggersCountFromBase
}

func (r *RespinGame) GetReTriggersCountFromFree() int {
	return r.ReTriggersCountFromFree
}

func (r *RespinGame) Merge(other *RespinGame) {
	r.WinsFromBase += other.WinsFromBase
	r.WinsFromFree += other.WinsFromFree
	r.ReTriggersCountFromBase += other.ReTriggersCountFromBase
	r.ReTriggersCountFromFree += other.ReTriggersCountFromFree
	r.ReTriggerStepsCountFromBase += other.ReTriggerStepsCountFromBase
	r.ReTriggerStepsCountFromFree += other.ReTriggerStepsCountFromFree
}
