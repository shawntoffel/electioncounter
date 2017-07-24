package electioncounter

import (
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
	Count(config Config) (*counters.Result, error)
}

type electionCounter struct {
}

func NewElectionCounter() ElectionCounter {
	return &electionCounter{}
}

func (c *electionCounter) Count(config Config) (*counters.Result, error) {
	counterFactory := NewCounterFactory()
	counter, err := counterFactory.GetCounter(config.Method)

	if err != nil {
		return nil, err
	}

	counterConfig := counters.CounterConfig{}
	counterConfig.NumberToElect = config.NumberToElect
	counterConfig.Ballots = config.Ballots
	counterConfig.Candidates = config.Candidates
	counterConfig.Precision = config.Precision

	counter.Create(counterConfig)

	for {
		if counter.HasEnded() {
			break
		}

		counter.UpdateRound()
	}

	result := counter.Result()
	return &result, nil
}
