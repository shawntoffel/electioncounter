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
	Candidates    []Candidate
	Precision     int
}

type Candidate struct {
	Id   string
	Name string
}

type Ballots []*list.List

type Events []Event
type Event struct {
	Description string
}

type Result struct {
	Candidates []Candidate
	Events     []Event
}

type counter struct {
	Changes         []Event
	ExpectedVersion int
}
