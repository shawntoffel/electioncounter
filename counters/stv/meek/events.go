package meek

import (
	"bytes"
	"fmt"
	"github.com/shawntoffel/electioncounter/counters"
)

type MeekEvent interface {
	Transition(m *meekStvCounter) string
}

type CreateCount struct {
	NumberToElect int
	Ballots       counters.Ballots
	Candidates    counters.Candidates
	Precision     int
}

func (e *CreateCount) Transition(state *meekStvCounter) string {
	state.NumberToElect = e.NumberToElect
	state.Precision = e.Precision

	state.Pool.AddNewCandidates(e.Candidates)
	state.Ballots = counters.Rollup(e.Ballots)

	buffer := bytes.Buffer{}

	buffer.WriteString("A new Meek STV count has been created")
	buffer.WriteString(fmt.Sprintf("\nCandidates: %d", len(e.Candidates)))
	buffer.WriteString(fmt.Sprintf("\nBallots: %d", len(e.Ballots)))
	buffer.WriteString(fmt.Sprintf("\nWinners: %d", state.NumberToElect))
	buffer.WriteString(fmt.Sprintf("\nPrecision: %d", state.Precision))

	return buffer.String()
}
