package main

import (
	"fmt"
)

type Event struct {
	Description string
}

type CounterEvent interface {
	Transition(c *counter)
}

type EventTask func(c *counter)

type CountCreated struct {
	Event

	NumberToElect int
	Ballots       []Ballot
	Candidates    []Candidate
	Precision     int
}

func (e CountCreated) Transition(c *counter) {
	c.NumberToElect = e.NumberToElect
	c.Ballots = e.Ballots

	c.Precision = e.Precision
	c.Pool = NewPool(e.Candidates)

	c.Description = "A new count has started."
}

type QuotaUpdated struct {
	Event
	NewQuota int64
}

func (e QuotaUpdated) Transition(c *counter) {
	c.Quota = e.NewQuota
	e.Description = fmt.Sprintf("Quota has been updated to %d", e.NewQuota)
}

type RoundEnded struct {
	Event
}

func (e RoundEnded) Transition(c *counter) {
	e.Description = fmt.Sprintf("Round %d has ended.", c.Round)

	c.Round++
}

type CandidateKeepValueUpdated struct {
	Event
	Id           string
	NewKeepValue int64
}

func (e CandidateKeepValueUpdated) Transition(c *counter) {
	c.Pool.SetKeepValue(e.Id, e.NewKeepValue)

	e.Description = fmt.Sprintf("The keep value for candidate '%s' has been updated to %d", e.Id, e.NewKeepValue)
}

type CandidateVotesUpdated struct {
	Event
	Id       string
	NewVotes int64
}

func (e CandidateVotesUpdated) Transition(c *counter) {
	c.Pool.SetVotes(e.Id, e.NewVotes)

	e.Description = fmt.Sprintf("The vote count for candidate '%s' has been updated to %d", e.Id, e.NewVotes)
}

type ElectCandidate struct {
	Event
	Id string
}

func (e ElectCandidate) Transition(c *counter) {
	c.Pool.SetStatus(e.Id, Elected)

	e.Description = fmt.Sprintf("Candidate '%s' has been elected.", e.Id)
}

type ExcludeCandidate struct {
	Event
	Id string
}

func (e ExcludeCandidate) Transition(c *counter) {
	c.Pool.SetStatus(e.Id, Excluded)

	e.Description = fmt.Sprintf("Candidate '%s' has been excluded.", e.Id)
}
