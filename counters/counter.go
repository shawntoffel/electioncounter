package counters

import (
	"container/list"
)

type Counter interface {
	Create(counterConfig CounterConfig)
	UpdateRound()
	HasEnded() bool
	Result() Result
}

type CounterConfig struct {
	NumberToElect int
	Ballots       Ballots
	Candidates    Candidates
	Precision     int
}

type Candidates []Candidate
type Candidate struct {
	Id   string
	Name string
}

type Ballot *list.List
type Ballots []*list.List

type Events []Event
type Event struct {
	Description string
}

type Result struct {
	Candidates Candidates
	Events     Events
}

type counter struct {
	Changes         []Event
	ExpectedVersion int
}
