package main

type MeekStv interface {
	Stv
	Initialize(config StvConfig)

	Run() ([]Candidate, Events)
}

type meekStv struct {
	stv

	Scaler int
}

func NewMeekStv() MeekStv {
	m := meekStv{}

	m.Counter = NewCounter(nil)

	return &m
}

func (m *meekStv) Initialize(config StvConfig) {
	m.Counter.Initialize(config)
}

func (m *meekStv) Run() ([]Candidate, Events) {
	m.Counter.SetInitialQuota()

	m.Counter.InitializeVotes()

	for {
		if m.Counter.HasEnded() {
			break
		}

		m.Counter.UpdateRound()
	}

	return m.Counter.Results()
}
