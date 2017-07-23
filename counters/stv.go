package counters

type StvCounter interface {
	CreateCount(config StvConfig)
	SetInitialQuota()
	InitializeVotes()
	UpdateRound()
	HasEnded() bool
	Results() ([]Result, Events)
	Events() []Event
}

type stvCounter struct {
	Changes         Events
	ExpectedVersion int
}

type StvConfig struct {
	NumberToElect int
	Ballots       Ballots
	Candidates    []Candidate
	Precision     int
}

type Result struct {
	CandidateId string
	Votes       int64
}

func GetScaler(precision int) int64 {
	var scaler = int64(1)
	for i := 0; i < precision; i++ {
		scaler *= 10
	}

	return scaler
}
