package events

import (
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
)

type ExcludeRemaining struct{}

func (e *ExcludeRemaining) Transition(s *state.MeekState) (string, error) {
	candidates := s.Pool.Candidates()

	for _, candidate := range candidates {
		if candidate.Status != state.Elected {
			s.Pool.Exclude(candidate.Id)
		}
	}

	return "All remaining candidates have been excluded.", nil
}
