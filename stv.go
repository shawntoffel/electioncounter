package electioncounter

import (
	"github.com/shawntoffel/electioncounter/counters"
)

type Stv interface {
	Initialize(config counters.StvConfig)
	Run() ([]counters.Candidate, counters.Events)
}

type stv struct {
	StvCounter counters.StvCounter
}

func NewStv(stvCounter counters.StvCounter) Stv {
	s := stv{}

	s.StvCounter = stvCounter

	return &s
}

func (s *stv) Initialize(config counters.StvConfig) {
	s.StvCounter.Initialize(config)
}

func (s *stv) Run() ([]counters.Candidate, counters.Events) {
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
