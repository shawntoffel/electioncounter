package events

import (
	"bytes"
	"fmt"
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
	"github.com/shawntoffel/electioncounter/election"
	"github.com/shawntoffel/math"
)

type Create struct {
	NumSeats   int
	Ballots    election.Ballots
	Candidates election.Candidates
	Precision  int
}

func (e *Create) Transition(state *state.MeekState) (string, error) {
	state.NumSeats = e.NumSeats
	state.Precision = e.Precision
	state.Pool.AddNewCandidates(e.Candidates)
	state.Ballots = e.Ballots.Rollup()

	state.Scale = math.Pow64(int64(10), int64(state.Precision))

	buffer := bytes.Buffer{}

	buffer.WriteString("A new Meek STV count has been created")
	buffer.WriteString(fmt.Sprintf("\nCandidates: %d", len(e.Candidates)))
	buffer.WriteString(fmt.Sprintf("\nBallots: %d", len(e.Ballots)))
	buffer.WriteString(fmt.Sprintf("\nSeats: %d", state.NumSeats))
	buffer.WriteString(fmt.Sprintf("\nPrecision: %d", state.Precision))

	return buffer.String(), nil
}
