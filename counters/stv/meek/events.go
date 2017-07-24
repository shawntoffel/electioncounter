package meek

import (
	"bytes"
	"fmt"
	"github.com/shawntoffel/electioncounter/counters"
)

type MeekEvent interface {
	Transition(m *meekStvCounter)
	Describe() string
}

type meekEvent struct {
	counters.Event
}

func (e *meekEvent) Describe() string {
	return e.Description
}

type CreateCount struct {
	meekEvent

	NumberToElect int
	Ballots       counters.Ballots
	Candidates    counters.Candidates
	Precision     int
}

func (e *CreateCount) Transition(state *meekStvCounter) {
	state.NumberToElect = e.NumberToElect
	state.Precision = e.Precision

	state.Pool.AddNewCandidates(e.Candidates)
	state.Ballots = counters.Rollup(e.Ballots)

	buffer := bytes.Buffer{}

	buffer.WriteString("A new Meek STV count has been created")
	buffer.WriteString(fmt.Sprintf("\nCandidates: %d", len(e.Candidates)))
	buffer.WriteString(fmt.Sprintf("\nBallots: %d", len(e.Ballots)))
	buffer.WriteString(fmt.Sprintf("\nPrecision: %d", state.Precision))
	buffer.WriteString(fmt.Sprintf("\nHow many to elect: %d", state.NumberToElect))

	e.Description = buffer.String()
}
