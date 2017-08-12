package event

import (
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
)

type MeekEvent interface {
	Process() (election.Event, error)
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
	Candidates state.Candidates
}
