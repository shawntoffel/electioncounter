package event

import (
	"bytes"
	"fmt"
	"github.com/shawntoffel/electioncounter/election"
)

func (e *CountCreated) Process() election.Event {
	buffer := bytes.Buffer{}

	buffer.WriteString("A new Meek STV count has been created")
	buffer.WriteString(fmt.Sprintf("\nCandidates: %d", len(e.Candidates)))
	buffer.WriteString(fmt.Sprintf("\nBallots: %d", len(e.Ballots)))
	buffer.WriteString(fmt.Sprintf("\nSeats: %d", e.NumSeats))
	buffer.WriteString(fmt.Sprintf("\nPrecision: %d", e.Precision))

	description := buffer.String()

	return election.Event{description}
}

func (e *CandidatesExcluded) Process() election.Event {
	names := []string{}

	for _, candidate := range e.Candidates {
		names = append(names, candidate.Name)
	}

	description := fmt.Sprintf("The following candidates have been excluded: %v", names)

	return election.Event{description}
}

func (e *AllHopefulCandidatesElected) Process() election.Event {
	description := "All hopeful candidates have been elected."

	return election.Event{description}
}

func (e *RoundStarted) Process() election.Event {
	description := fmt.Sprintf("Round %d has started.", e.Round)

	return election.Event{description}
}

func (e *LowestCandidateExcluded) Process() election.Event {

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

	return election.Event{description}
}

func (e *RemainingCandidatesExcluded) Process() election.Event {
	description := "All remaining candidates have been excluded."

	return election.Event{description}
}
