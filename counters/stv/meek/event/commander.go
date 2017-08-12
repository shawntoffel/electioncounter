package event

import (
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
	"github.com/shawntoffel/electioncounter/election"
)

type Commander interface {
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

type commander struct {
	ElectedAll bool
	State      state.MeekState
	Consumer   Consumer
}

func NewCommander() Commander {
	c := commander{}
	c.State = state.NewMeekState()
	c.Consumer = NewConsumer()

	return &c
}

func (c *commander) Create(config election.Config) {
	c.State.NumSeats = e.NumSeats
	c.State.Precision = e.Precision
	c.State.Pool.AddNewCandidates(e.Candidates)
	c.State.Ballots = e.Ballots.Rollup()

	state.Scale = math.Pow64(10, int64(state.Precision))
	state.MaxIterations = 1000

	event := CountCreated{}
	event.Candidates = config.Candidates
	event.Ballots = config.Ballots
	event.Precision = config.Precision
	event.NumSeats = config.NumSeats

	p.Consumer.ProcessEvent(&event)
}

func (c *commander) ExcludeWithdrawnCandidates(ids []string) {
	excluded := state.MeekCandidates{}

	for _, id := range ids {
		candidate := state.Pool.Exclude(id)

		excluded = append(excluded, candidate)
	}

	p.Consumer.ProcessEvent(&CandidatesExcluded{excluded})
}

func (c *commander) PerformPreliminaryElection() {
	numCandidates := p.State.Pool.Count()
	numExcluded := p.State.Pool.ExcludedCount()

	if numCandidates <= (p.State.NumSeats + numExcluded) {
		c.State.Pool.ElectHopeful()
		s.State.ElectedAll = true
		p.Consumer.ProcessEvent(&AllHopefulCandidatesElected{})
	}
}

func (c *commander) IncrementRound() {
	c.State.Round = c.State.Round + 1

	p.Consumer.ProcessEvent(&RoundStarted{c.State.Round})
}

func (c *commander) ExcludeLowestCandidate() {
	lowestCandidates := c.Pool.Lowest()

	toExclude := lowestCandidates[0]

	randomUsed := false

	if len(lowestCandidates) > 1 {
		seed := rand.NewSource(time.Now().Unix())
		r := rand.New(seed)
		i := r.Intn(len(lowestCandidates))
		toExclude = lowestCandidates[i]

		randomUsed = true
	}

	c.Pool.Exclude(toExclude.Id)
	c.Consumer.ProcessEvent(&LowestCandidateExcluded{lowestCandidates, toExclude, randomUsed})
}

func (c *commander) ExcludeRemainingCandidates() {
	candidates := c.Pool.Candidates()

	for _, candidate := range candidates {
		if candidate.Status != state.Elected {
			s.Pool.Exclude(candidate.Id)
		}
	}

	p.Consumer.ProcessEvent(&RemainingCandidatesExcluded{})
}

func (p *commander) DistributeVotes() {

	/*
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

		}*/
}

func (c *commander) HasEnded() bool {
	if c.Error != nil {
		return true
	}

	if c.State.ElectedAll {
		return true
	}

	numElected := c.State.Pool.ElectedCount()

	return numElected == c.State.NumSeats
}

func (c *commander) RoundHasEnded() bool {

	if !c.State.MeekRound.AnyElected {
		return true
	}

	numElected := c.State.Pool.ElectedCount()

	if numElected >= c.State.NumSeats {
		return true
	}

	return false
}

func (c *commander) Changes() (election.Events, error) {
	//return c.State.Events, c.Error

	return election.Events{}, c.Error
}
