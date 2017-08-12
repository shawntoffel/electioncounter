package event

import (
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
	"github.com/shawntoffel/electioncounter/election"
)

type MeekEvent interface {
	Process() election.Event
}

type CountCreated struct {
	NumSeats   int
	Ballots    election.Ballots
	Candidates election.Candidates
	Precision  int
}

type CandidatesExcluded struct {
	Candidates state.MeekCandidates
}

type AllHopefulCandidatesElected struct{}

type RoundStarted struct {
	Round int
}

type LowestCandidateExcluded struct {
	LowestCandidates  state.MeekCandidates
	ExcludedCandidate state.MeekCandidate
	RandomUsed        bool
}

type RemainingCandidatesExcluded struct {
	Candidates state.MeekCandidates
}
