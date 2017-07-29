package meek

import (
	"github.com/shawntoffel/electioncounter/election"
)

type CandidateStatus int

const (
	Elected  CandidateStatus = iota
	Hopeful  CandidateStatus = iota
	Excluded CandidateStatus = iota
)

type MeekCandidate struct {
	election.Candidate
	Status CandidateStatus
	Weight int64
	Votes  int64
}
