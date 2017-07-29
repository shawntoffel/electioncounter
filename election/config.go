package election

type Config struct {
	NumberToElect       int
	Ballots             Ballots
	Candidates          Candidates
	WithdrawnCandidates []string
	Precision           int
}
