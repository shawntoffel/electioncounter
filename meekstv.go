package main

type MeekStv interface {
	Stv
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

	countCreated := CountCreated{}
	countCreated.Candidates = config.Candidates
	countCreated.Ballots = config.Ballots
	countCreated.Precision = config.Precision
	countCreated.NumberToElect = config.NumberToElect

	m.Counter.HandleEvent(countCreated)
}

func (m *meekStv) Run() {
	m.Counter.SetInitialQuota()

	m.Counter.InitializeVotes()

	for {
		if m.Counter.HasEnded() {
			break
		}

		m.Counter.UpdateRound()
	}

}
