package counters

import (
	"container/list"
)

type Events []CounterEvent

type StvCounter interface {
	HandleEvent(event CounterEvent)
	Events() Events

	Initialize(config StvConfig)

	SetInitialQuota()
	InitializeVotes()
	UpdateRound()
	HasEnded() bool
	Results() ([]Candidate, Events)
}

type stvCounter struct {
	Changes         Events
	ExpectedVersion int
}

type Ballots []*list.List

type StvConfig struct {
	NumberToElect int
	Ballots       Ballots
	Candidates    []Candidate
	Precision     int
}

func GetScaler(precision int) int64 {
	var scaler = int64(1)
	for i := 0; i < precision; i++ {
		scaler *= 10
	}

	return scaler
}
