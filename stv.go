package electioncounter

import (
	"github.com/shawntoffel/electioncounter/counters"
)

type Stv interface {
	Run(config counters.StvConfig) ([]counters.Candidate, counters.Events)
}

type stv struct {
	StvCounter counters.StvCounter
}

func NewStv(stvCounter counters.StvCounter) Stv {
	return &stv{stvCounter}
}

func (s *stv) Run(config counters.StvConfig) ([]counters.Candidate, counters.Events) {
	s.StvCounter.CreateCount(config)
	s.StvCounter.SetInitialQuota()

	s.StvCounter.InitializeVotes()

	for {
		if s.StvCounter.HasEnded() {
			break
		}

		s.StvCounter.UpdateRound()
	}

	return s.StvCounter.Results()
}
