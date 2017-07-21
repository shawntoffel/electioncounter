package main

import (
	"container/list"
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
	Ballots       []*list.List
	Candidates    []Candidate
	Precision     int
}

func (e CountCreated) Transition(c *counter) {
	c.NumberToElect = e.NumberToElect
	c.Ballots = e.Ballots

	c.Precision = e.Precision
	c.Pool = NewPool(e.Candidates)

	c.Scaler = GetScaler(c.Precision)

	e.Description = "A new count has started."
}

type QuotaUpdated struct {
	Event
	NewQuota int64
}

func (e QuotaUpdated) Transition(c *counter) {
	c.Quota = e.NewQuota
	e.Description = fmt.Sprintf("Quota has been updated to %d", e.NewQuota)
}

type SetDroopQuota struct {
	Event
}

func (e SetDroopQuota) Transition(c *counter) {

	var numBallots = int64(len(c.Ballots))

	var droop = ((numBallots*c.Scaler)/(int64(c.NumberToElect)+1))/c.Scaler*c.Scaler + c.Scaler

	c.Quota = droop

	e.Description = fmt.Sprintf("Quota has been set via Droop to %d", c.Quota)
}

type IncrementRound struct {
	Event
}

func (e IncrementRound) Transition(c *counter) {
	c.Round++
	e.Description = fmt.Sprintf("Round %d has started.", c.Round)
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
	Id string
	Event
}

func (e ElectCandidate) Transition(c *counter) {
	c.Pool.SetStatus(e.Id, Elected)

	e.Description = fmt.Sprintf("Candidate '%s' has been elected.", e.Id)
}

type ExcludeCandidate struct {
	Id string
	Event
}

func (e ExcludeCandidate) Transition(c *counter) {
	c.Pool.SetStatus(e.Id, Excluded)

	e.Description = fmt.Sprintf("Candidate '%s' has been excluded.", e.Id)
}

type InitializeVotes struct {
	Event
}

func (e InitializeVotes) Transition(c *counter) {

}
