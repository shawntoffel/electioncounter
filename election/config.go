package election

type Config struct {
	NumSeats            int
	Ballots             Ballots
	Candidates          Candidates
	WithdrawnCandidates []string
	Precision           int
}
