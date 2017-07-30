package state

import (
	"github.com/shawntoffel/electioncounter/election"
)

type CandidateStatus int

const (
	Elected  CandidateStatus = iota
	Hopeful  CandidateStatus = iota
	Excluded CandidateStatus = iota
)

type MeekCandidates []MeekCandidate
type MeekCandidate struct {
	election.Candidate
	Status CandidateStatus
	Weight int64
	Votes  int64
}

type ByVotes MeekCandidates

func (c ByVotes) Len() int {
	return len(c)
}

func (c ByVotes) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ByVotes) Less(i, j int) bool {
	return c[i].Votes < c[j].Votes
}
