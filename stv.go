package main

type Stv interface {
	Initialize(config StvConfig)
	Run() ([]Candidate, Events)
}

type stv struct {
	StvCounter StvCounter
}

func NewStv(stvCounter StvCounter) Stv {
	s := stv{}

	s.StvCounter = stvCounter

	return &s
}

func (s *stv) Initialize(config StvConfig) {
	s.StvCounter.Initialize(config)
}

func (s *stv) Run() ([]Candidate, Events) {
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
