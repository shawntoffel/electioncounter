package electioncounter

import (
	"github.com/shawntoffel/electioncounter/election"
	"github.com/shawntoffel/electioncounter/factory"
)

type ElectionCounter interface {
	Count(method string, config election.Config) (*election.Result, error)
}

type electionCounter struct{}

func NewElectionCounter() ElectionCounter {
	return &electionCounter{}
}

func (c *electionCounter) Count(method string, config election.Config) (*election.Result, error) {
	counter, err := factory.NewCounter(method)

	if err != nil {
		return nil, err
	}

	counter.Initialize(config)

	for {
		if counter.HasEnded() {
			break
		}

		counter.UpdateRound()
	}

	return counter.Result()
}
