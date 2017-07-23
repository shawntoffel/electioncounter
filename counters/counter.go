package counters

type Ballots []*list.List
type Counter interface {
	CreateCount(countConfig CountConfig)
	UpdateRound()
	HasEnded() bool
	Results() []Result
}

type CountConfig struct {
}

type Result struct {
	Candidates []Candidate
	Events     []Event
}

type Candidate struct {
	Id   string
	Name string
}

type Event struct {
	Description string
}

type counter struct {
	Changes         Events
	ExpectedVersion int
}
