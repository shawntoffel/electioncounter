package state

import (
	"github.com/shawntoffel/electioncounter/counters/stv"
)

type MeekState struct {
	stv.State
	MeekRound     MeekRound
	Pool          Pool
	Precision     int
	Scale         int64
	ElectedAll    bool
	MaxIterations int
}

func NewMeekState() *MeekState {
	meekState := MeekState{}
	meekState.Pool = NewPool()
	meekState.MeekRound = MeekRound{}

	return &meekState
}

type MeekRound struct {
	Excess     int64
	AnyElected bool
}
