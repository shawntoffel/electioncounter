package event

import (
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
)

type MeekEvent interface {
	Transition(m *state.MeekState) (string, error)
}
