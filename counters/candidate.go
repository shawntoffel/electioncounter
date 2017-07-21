package counters

type CandidateStatus int

const (
	Elected  CandidateStatus = iota
	Hopeful  CandidateStatus = iota
	Excluded CandidateStatus = iota
)

type Candidate struct {
	Id             string
	Status         CandidateStatus
	KeepValue      int64
	Votes          int64
	FirstRankCount int
}
