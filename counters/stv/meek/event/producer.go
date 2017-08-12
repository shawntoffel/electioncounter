package event

import (
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
	"github.com/shawntoffel/electioncounter/election"
)

type Producer interface {
	Create(config election.Config)
	ExcludeWithdrawnCandidates(ids []string)
	IncrementRound()
	PerformPreliminaryElection()
	ExcludeRemainingCandidates()
	ExcludeLowestCandidate()
	HasEnded() bool
	DistributeVotes()
	RoundHasEnded() bool
	Changes() (election.Events, error)
}

type producer struct {
	ElectedAll bool
	Consumer   Consumer
}

func NewProducer(events []MeekEvent) MeekEventProducer {
	p := producer{}
	p.Consumer = NewConsumer()

	for _, event := range events {
		p.Consumer.ProcessEvent(event)
	}

	return &p
}

func (p *producer) Create(config election.Config) {
	event := Create{}
	event.Candidates = config.Candidates
	event.Ballots = config.Ballots
	event.Precision = config.Precision
	event.NumSeats = config.NumSeats

	p.Consumer.ProcessEvent(&event)
}

func (p *producer) ExcludeWithdrawnCandidates(ids []string) {
	p.Consumer.ProcessEvent(&WithdrawlCandidates{ids})
}

func (p *producer) PerformPreliminaryElection() {
	numCandidates := p.State.Pool.Count()
	numExcluded := p.State.Pool.ExcludedCount()

	if numCandidates <= (p.State.NumSeats + numExcluded) {
		p.Consumer.ProcessEvent(&ElectAll{})
	}
}

func (p *producer) HasEnded() bool {
	if p.Error != nil {
		return true
	}

	if p.State.ElectedAll {
		return true
	}

	numElected := p.State.Pool.ElectedCount()

	return numElected == p.State.NumSeats
}

func (p *producer) ExcludeRemainingCandidates() {
	p.Consumer.ProcessEvent(&ExcludeRemaining{})
}

func (p *producer) ExcludeLowestCandidate() {
	p.Consumer.ProcessEvent(&ExcludeLowest{})
}

func (p *producer) IncrementRound() {
	p.Consumer.ProcessEvent(&IncrementRound{})
}

func (p *producer) DistributeVotes() {

	for i := 0; i < p.State.MaxIterations; i++ {
		p.State.MeekRound.Excess = 0

		for _, ballot := range p.State.Ballots {
			value := int64(ballot.Count) * p.State.Scale

			ended := false

			iter := ballot.Ballot.List.Front()

			for {
				candidate := p.State.Pool.Candidate(iter.Value.(string))

				if !ended && candidate.Weight > 0 {
					ended = candidate.Status == state.Hopeful

					if ended {
						votes := candidate.Votes + value
						p.State.Pool.SetVotes(candidate.Id, votes)
						value = 0
					} else {
						votes := candidate.Votes + value*candidate.Weight
						p.State.Pool.SetVotes(candidate.Id, votes)
						value = value * (p.State.Scale - candidate.Weight) / p.State.Scale
					}
				}

				if iter.Next() == nil {
					break
				}

				iter = iter.Next()
			}

			p.State.MeekRound.Excess = p.State.MeekRound.Excess + value

		}

		p.State.Quota = 1600000

		break

	}
}

func (p *producer) RoundHasEnded() bool {

	if !p.State.MeekRound.AnyElected {
		return true
	}

	numElected := p.State.Pool.ElectedCount()

	if numElected >= p.State.NumSeats {
		return true
	}

	return false
}

func (p *producer) Changes() (election.Events, error) {
	return p.State.Events, p.Error
}
