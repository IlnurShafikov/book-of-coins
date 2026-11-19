package result

// FreeGame accumulates statistics for the free game phase across a simulation run.
type FreeGame struct {
	// Wins is the sum of all win amounts from free game spins.
	Wins uint64
	// RespinGameTriggered is the number of respin sessions triggered by free spins.
	RespinGameTriggered uint64
	// AllRespinGameSteps is the total number of respins played after free spins.
	AllRespinGameSteps int
}

// newFreeGame creates a FreeGame with zeroed statistics.
func newFreeGame() *FreeGame {
	return &FreeGame{
		Wins:                0,
		RespinGameTriggered: 0,
		AllRespinGameSteps:  0,
	}
}

// AddWins adds the given win amount to the free game win total.
func (f *FreeGame) AddWins(wins uint64) {
	f.Wins += wins
}

// AddRespinGameSteps records the respins played after a free spin.
// Sessions without respins (steps <= 0) are ignored.
//
// Parameters:
//   - steps: number of respins played
func (f *FreeGame) AddRespinGameSteps(steps int) {
	if steps > 0 {
		f.AllRespinGameSteps += steps
		f.RespinGameTriggered++
	}
}

// GetWins returns the total win amount accumulated from free game spins.
func (f *FreeGame) GetWins() uint64 {
	return f.Wins
}

// GetRespinGameTriggered returns the number of respin sessions triggered by free spins.
func (f *FreeGame) GetRespinGameTriggered() uint64 {
	return f.RespinGameTriggered
}

// GetAllRespinGameSteps returns the total number of respins played after free spins.
func (f *FreeGame) GetAllRespinGameSteps() int {
	return f.AllRespinGameSteps
}

// Merge adds another FreeGame's statistics into this one.
func (f *FreeGame) Merge(other *FreeGame) {
	f.Wins += other.Wins
	f.RespinGameTriggered += other.RespinGameTriggered
	f.AllRespinGameSteps += other.AllRespinGameSteps
}
