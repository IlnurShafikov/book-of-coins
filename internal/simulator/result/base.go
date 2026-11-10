package result

type BaseGame struct {
	Wins                uint64
	FreeGameTriggered   uint64
	AllFreeGameSteps    int
	RespinGameTriggered uint64
	AllRespinGameSteps  int
}

func newBaseGame() *BaseGame {
	return &BaseGame{
		Wins:                0,
		FreeGameTriggered:   0,
		AllFreeGameSteps:    0,
		RespinGameTriggered: 0,
		AllRespinGameSteps:  0,
	}
}

func (b *BaseGame) AddWins(Wins uint64) {
	b.Wins += Wins
}

func (b *BaseGame) AddFreeGameSteps(steps int) {
	if steps > 0 {
		b.AllFreeGameSteps += steps
		b.FreeGameTriggered++
	}
}

func (b *BaseGame) AddRespinGameSteps(steps int) {
	if steps > 0 {
		b.AllRespinGameSteps += steps
		b.RespinGameTriggered++
	}
}

func (b *BaseGame) GetFreeGameTriggered() uint64 {
	return b.FreeGameTriggered
}

func (b *BaseGame) GetRespinGameTriggered() uint64 {
	return b.RespinGameTriggered
}

func (b *BaseGame) GetAllFreeGameSteps() int {
	return b.AllFreeGameSteps
}

func (b *BaseGame) GetAllRespinGameSteps() int {
	return b.AllRespinGameSteps
}

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
