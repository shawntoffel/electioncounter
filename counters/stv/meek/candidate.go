package meek

import (
	"github.com/shawntoffel/electioncounter/counters"
)

type CandidateStatus int

const (
	Elected  CandidateStatus = iota
	Hopeful  CandidateStatus = iota
	Excluded CandidateStatus = iota
)

type MeekCandidate struct {
	counters.Candidate
	Status CandidateStatus
	Weight int64
	Votes  int64
}
