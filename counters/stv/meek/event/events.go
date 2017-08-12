package event

type MeekEvent interface {
	Process(m *state.MeekState) (string, error)
}

type Create struct {
	NumSeats   int
	Ballots    election.Ballots
	Candidates election.Candidates
	Precision  int
}

type ElectAll struct{}

type ExcludeLowest struct{}

type ExcludeRemaining struct{}

type IncrementRound struct{}

type WithdrawlCandidates struct {
	Ids []string
}
