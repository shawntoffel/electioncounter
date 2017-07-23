package electioncounter

import (
	"container/list"

	"github.com/shawntoffel/electioncounter/counters"
)

type Config struct {
	NumberToElect int
	Ballots       counters.Ballots
	Candidates    []counters.Candidate
	Precision     int
	Method        string
}

type ElectionCounter interface {
	Count(config Config) []Result
}

type electionCounter struct {
}

func NewElectionCounter() ElectionCounter {
	return &electionCounter{}
}

func (c *electionCounter) Count(config Config) {
	counterFactory := counters.NewCounterFactory()
	counter := counterFactory.GetCounter(config.Name)

	for {
		if counter.HasEnded() {
			break
		}

		counter.UpdateRound()
	}

	return counter.Results()
}
