package events

import (
	"fmt"
)

type WithdrawlCandidates struct {
	Ids []string
}

func (e *WithdrawlCandidates) Transition(state *meekState) string {
	names := []string{}

	for _, id := range e.Ids {
		c := state.Pool.Candidate(id)
		names = append(names, c.Name)
		state.Pool.Exclude(id)
	}

	return fmt.Sprintf("The following candidates have been excluded: %v", names)
}
