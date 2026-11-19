package result

// BaseGame accumulates statistics for the base game phase across a simulation run.
type BaseGame struct {
	// Wins is the sum of all win amounts from base game spins.
	Wins uint64
	// FreeGameTriggered is the number of rounds that triggered the free game.
	FreeGameTriggered uint64
	// AllFreeGameSteps is the total number of free spins awarded by the base game.
	AllFreeGameSteps int
	// RespinGameTriggered is the number of rounds in which the base game triggered respins.
	RespinGameTriggered uint64
	// AllRespinGameSteps is the total number of respins played after the base game.
	AllRespinGameSteps int
}

// newBaseGame creates a BaseGame with zeroed statistics.
func newBaseGame() *BaseGame {
	return &BaseGame{
		Wins:                0,
		FreeGameTriggered:   0,
		AllFreeGameSteps:    0,
		RespinGameTriggered: 0,
		AllRespinGameSteps:  0,
	}
}

// AddWins adds the given win amount to the base game win total.
func (b *BaseGame) AddWins(wins uint64) {
	b.Wins += wins
}

// AddFreeGameSteps records the free spins awarded by a base spin.
// Spins without free spins (steps <= 0) are ignored.
//
// Parameters:
//   - steps: number of free spins awarded
func (b *BaseGame) AddFreeGameSteps(steps int) {
	if steps > 0 {
		b.AllFreeGameSteps += steps
		b.FreeGameTriggered++
	}
}

// AddRespinGameSteps records the respins played after a base spin.
// Sessions without respins (steps <= 0) are ignored.
//
// Parameters:
//   - steps: number of respins played
func (b *BaseGame) AddRespinGameSteps(steps int) {
	if steps > 0 {
		b.AllRespinGameSteps += steps
		b.RespinGameTriggered++
	}
}

// GetFreeGameTriggered returns the number of rounds that triggered the free game.
func (b *BaseGame) GetFreeGameTriggered() uint64 {
	return b.FreeGameTriggered
}

// GetRespinGameTriggered returns the number of rounds in which the base game triggered respins.
func (b *BaseGame) GetRespinGameTriggered() uint64 {
	return b.RespinGameTriggered
}

// GetAllFreeGameSteps returns the total number of free spins awarded by the base game.
func (b *BaseGame) GetAllFreeGameSteps() int {
	return b.AllFreeGameSteps
}

// GetAllRespinGameSteps returns the total number of respins played after the base game.
func (b *BaseGame) GetAllRespinGameSteps() int {
	return b.AllRespinGameSteps
}

// GetWins returns the total win amount accumulated from base game spins.
func (b *BaseGame) GetWins() uint64 {
	return b.Wins
}

// Merge adds another BaseGame's statistics into this one.
func (b *BaseGame) Merge(other *BaseGame) {
	b.Wins += other.Wins
	b.AllFreeGameSteps += other.AllFreeGameSteps
	b.AllRespinGameSteps += other.AllRespinGameSteps
	b.FreeGameTriggered += other.FreeGameTriggered
	b.RespinGameTriggered += other.RespinGameTriggered
}
