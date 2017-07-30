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
