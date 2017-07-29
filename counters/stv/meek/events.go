package meek

import (
	"bytes"
	"fmt"
	"github.com/shawntoffel/electioncounter/election"
	"github.com/shawntoffel/math"
)

type MeekEvent interface {
	Transition(m *meekStv) string
}

type CreateCount struct {
	NumberToElect int
	Ballots       election.Ballots
	Candidates    election.Candidates
	Precision     int
}

func (e *CreateCount) Transition(state *meekStv) string {
	state.NumberToElect = e.NumberToElect
	state.Precision = e.Precision
	state.Pool.AddNewCandidates(e.Candidates)
	state.Ballots = e.Ballots.Rollup()

	state.Scale = int64(math.Pow(10, state.Precision))

	buffer := bytes.Buffer{}

	buffer.WriteString("A new Meek STV count has been created")
	buffer.WriteString(fmt.Sprintf("\nCandidates: %d", len(e.Candidates)))
	buffer.WriteString(fmt.Sprintf("\nBallots: %d", len(e.Ballots)))
	buffer.WriteString(fmt.Sprintf("\nWinners: %d", state.NumberToElect))
	buffer.WriteString(fmt.Sprintf("\nPrecision: %d", state.Precision))

	return buffer.String()
}

type WithdrawlCandidates struct {
	Ids []string
}

func (e *WithdrawlCandidates) Transition(state *meekStv) string {
	names := []string{}

	for _, id := range e.Ids {
		c := state.Pool.Candidate(id)
		names = append(names, c.Name)
		state.Pool.Exclude(id)
	}

	return fmt.Sprintf("The following candidates have been excluded: %v", names)
}
