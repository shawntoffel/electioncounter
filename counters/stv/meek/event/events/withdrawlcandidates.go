package events

import (
	"fmt"
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
)

type WithdrawlCandidates struct {
	Ids []string
}

func (e *WithdrawlCandidates) Transition(state *state.MeekState) (string, error) {
	names := []string{}

	for _, id := range e.Ids {
		c := state.Pool.Candidate(id)
		names = append(names, c.Name)
		state.Pool.Exclude(id)
	}

	return fmt.Sprintf("The following candidates have been excluded: %v", names), nil
}
