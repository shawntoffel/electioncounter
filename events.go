package main

type CounterEvent interface {
	Task(c *counter)
}

type EventTask func(c *counter)

type CountCreated struct {
	NumberToElect int
	Ballots       []Ballot
	Candidates    []Candidate
	Precision     int
}

type QuotaUpdated struct {
	NewQuota int64
}

func (e QuotaUpdated) Task(c *counter) {
	c.Quota = e.NewQuota
}

type RoundEnded struct{}

type CandidateKeepValueUpdated struct {
	Id           string
	NewKeepValue int64
}

type CandidateVotesUpdated struct {
	Id       string
	NewVotes int64
}

type ElectCandidate struct {
	Id string
}

type ExcludeCandidate struct {
	Id string
}
