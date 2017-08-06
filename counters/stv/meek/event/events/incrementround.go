package events

import (
	"fmt"
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
)

type IncrementRound struct{}

func (e *IncrementRound) Transition(s *state.MeekState) (string, error) {
	s.Round = s.Round + 1

	return fmt.Sprintf("Round %d has started.", s.Round), nil
}
