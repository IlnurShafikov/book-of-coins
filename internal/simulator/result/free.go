package result

type FreeGame struct {
	Wins                uint64
	RespinGameTriggered uint64
	AllRespinGameSteps  int
}

func newFreeGame() *FreeGame {
	return &FreeGame{
		Wins:                0,
		RespinGameTriggered: 0,
		AllRespinGameSteps:  0,
	}
}

func (f *FreeGame) AddWins(wins uint64) {
	f.Wins += wins
}

func (f *FreeGame) AddRespinGameSteps(steps int) {
	if steps > 0 {
		f.AllRespinGameSteps += steps
		f.RespinGameTriggered++
	}
}

func (f *FreeGame) GetWins() uint64 {
	return f.Wins
}

func (f *FreeGame) GetRespinGameTriggered() uint64 {
	return f.RespinGameTriggered
}

func (f *FreeGame) GetAllRespinGameSteps() int {
	return f.AllRespinGameSteps
}

func (f *FreeGame) Merge(other *FreeGame) {
	f.Wins += other.Wins
	f.RespinGameTriggered += other.RespinGameTriggered
	f.AllRespinGameSteps += other.AllRespinGameSteps
}
