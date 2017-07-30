package events

import (
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
)

type ElectAll struct{}

func (e *ElectAll) Transition(s *state.MeekState) (string, error) {

	candidates := s.Pool.Candidates()

	for _, c := range candidates {
		if c.Status == state.Hopeful {
			s.Pool.ElectCandidate(c.Id)
		}
	}

	return "All hopeful candidates have been elected.", nil
}