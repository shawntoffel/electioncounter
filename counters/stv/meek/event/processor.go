package event

import (
	"bytes"
	"fmt"
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
	"github.com/shawntoffel/electioncounter/election"
	"github.com/shawntoffel/math"
	"math/rand"
	"time"
)

func (e *CountCreated) Process() (election.Event, error) {
	buffer := bytes.Buffer{}

	buffer.WriteString("A new Meek STV count has been created")
	buffer.WriteString(fmt.Sprintf("\nCandidates: %d", len(e.Candidates)))
	buffer.WriteString(fmt.Sprintf("\nBallots: %d", len(e.Ballots)))
	buffer.WriteString(fmt.Sprintf("\nSeats: %d", e.NumSeats))
	buffer.WriteString(fmt.Sprintf("\nPrecision: %d", e.state.Precision))

	description := buffer.String()

	return election.Event{description}, nil
}

func (e *CandidatesExcluded) Process() (election.Event, error) {
	names := []string{}

	for c, _ := range e.Candidates {
		names = append(names, c.Name)
	}

	description := fmt.Sprintf("The following candidates have been excluded: %v", names), nil

	return election.Event{description}, nil
}

func (e *AllHopefulCandidatesElected) Process() (election.Event, error) {
	description := "All hopeful candidates have been elected."

	return election.Event{description}, nil
}

func (e *RoundStarted) Process() (election.Event, error) {
	description := fmt.Sprintf("Round %d has started.", e.Round), nil

	return election.Event{description}, nil
}

func (e *LowestCandidateExcluded) Process() (election.Event, error) {

	buffer := bytes.Buffer{}

	if e.RandomUsed {
		names := []string{}

		for _, candidate := range e.LowestCandidates {
			names = append(names, candidate.Name)
		}

		buffer.WriteString(fmt.Sprintf("Candidates %v are tied for the lowest number of votes.", names))
		buffer.WriteString(fmt.Sprintf(" Candidate %s was randomly selected to break the tie.", e.ExcludedCandidate.Name))
	}

	buffer.WriteString(" Candidate '" + e.ExcludedCandidate.Name + "' has the lowest votes and has been excluded.")

	description := buffer.String()

	return election.Event{description}, nil
}

func (e *RemainingCandidatesExcluded) Process() (election.Event, error) {
	description := "All remaining candidates have been excluded."

	return election.Event{description}, nil
}
