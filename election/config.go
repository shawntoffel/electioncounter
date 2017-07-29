package election

type Config struct {
	NumberToElect       int
	Ballots             Ballots
	Candidates          Candidates
	WithdrawnCandidates Candidates
	Precision           int
}
