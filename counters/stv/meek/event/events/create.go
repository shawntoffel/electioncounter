package events

import (
	"bytes"
	"fmt"
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
	"github.com/shawntoffel/electioncounter/election"
	"github.com/shawntoffel/math"
)

type Create struct {
	NumberToElect int
	Ballots       election.Ballots
	Candidates    election.Candidates
	Precision     int
}

func (e *Create) Transition(state *state.MeekState) (string, error) {
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

	return buffer.String(), nil
}