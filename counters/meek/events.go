package counters

import (
	"container/list"
	"fmt"
)

type QuotaUpdated struct {
	Event
	NewQuota int64
}

func (e *QuotaUpdated) Transition(c *counter) {
	c.Quota = e.NewQuota
	e.EventDescription = fmt.Sprintf("Quota has been updated to %s", FormatScaledNumber(e.NewQuota, c.Scaler))
}

type IncrementRound struct {
	Event
}

func (e *IncrementRound) Transition(c *counter) {
	c.Round++
	e.EventDescription = fmt.Sprintf("Round %d has started", c.Round)
}

type CandidateKeepValueUpdated struct {
	Event
	Id           string
	NewKeepValue int64
}

func (e *CandidateKeepValueUpdated) Transition(c *counter) {
	c.Pool.SetKeepValue(e.Id, e.NewKeepValue)

	e.EventDescription = fmt.Sprintf("The keep value for candidate '%s' has been updated to %s", e.Id, FormatScaledNumber(e.NewKeepValue, c.Scaler))
}

type CandidateVotesUpdated struct {
	Event
	Id       string
	NewVotes int64
}

func (e *CandidateVotesUpdated) Transition(c *counter) {
	c.Pool.SetVotes(e.Id, e.NewVotes)

	e.EventDescription = fmt.Sprintf("The vote count for candidate '%s' has been updated to %s", e.Id, FormatScaledNumber(e.NewVotes, c.Scaler))
}

type CandidateVotesReceived struct {
	Event
	Id            string
	ReceivedVotes int64
}

func (e *CandidateVotesReceived) Transition(c *counter) {

	candidate := c.Pool.Candidate(e.Id)

	existingVotes := candidate.Votes

	updatedVotes := existingVotes + e.ReceivedVotes

	c.Pool.SetVotes(e.Id, updatedVotes)

	e.EventDescription = fmt.Sprintf("Candidate '%s' has received %s votes and now has %s votes", e.Id, FormatScaledNumber(e.ReceivedVotes, c.Scaler), FormatScaledNumber(updatedVotes, c.Scaler))
}

type ElectCandidate struct {
	Id string
	Event
}

func (e *ElectCandidate) Transition(c *counter) {
	c.Pool.SetStatus(e.Id, Elected)

	e.EventDescription = fmt.Sprintf("Candidate '%s' has been elected.", e.Id)
}

type ExcludeCandidate struct {
	Id string
	Event
}

func (e *ExcludeCandidate) Transition(c *counter) {
	c.Pool.SetStatus(e.Id, Excluded)

	e.EventDescription = fmt.Sprintf("Candidate '%s' has been excluded.", e.Id)
}

type DistributeVotes struct {
	Event
}

func (e *DistributeVotes) Transition(c *counter) {
	e.EventDescription = "Distributing votes"
}

func FormatScaledNumber(input int64, scale int64) string {
	var first = input / scale
	var second = input % scale

	return fmt.Sprintf("%d.%d", first, second)
}
