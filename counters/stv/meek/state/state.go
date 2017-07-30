package state

import (
	"github.com/shawntoffel/electioncounter/counters/stv"
)

type MeekState struct {
	stv.StvState
	Pool       Pool
	Precision  int
	Scale      int64
	ElectedAll bool
}

func NewMeekState() *MeekState {
	state := MeekState{}
	state.Round = stv.StvRound{0, false}
	state.Pool = NewPool()

	return &state
}
