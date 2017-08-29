package event

import (
	"fmt"
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
	"github.com/shawntoffel/electioncounter/election"
	"github.com/shawntoffel/math"
	"math/rand"
	"time"
)

type CommanderStatus struct {
	Error   error
	Events  election.Events
	Elected state.MeekCandidates
}

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
	Status() CommanderStatus
}

type commander struct {
	ElectedAll bool
	State      *state.MeekState
	Consumer   Consumer
	Error      error
}

func NewCommander() Commander {
	c := commander{}
	c.State = state.NewMeekState()
	c.Consumer = NewConsumer()

	return &c
}

func (c *commander) Create(config election.Config) {
	c.State.NumSeats = config.NumSeats
	c.State.Precision = config.Precision
	c.State.Pool.AddNewCandidates(config.Candidates)
	c.State.Ballots = config.Ballots.Rollup()

	c.State.Scale = math.Pow64(10, int64(c.State.Precision))
	c.State.MaxIterations = 1000

	event := CountCreated{}
	event.Candidates = config.Candidates
	event.Ballots = config.Ballots
	event.Precision = config.Precision
	event.NumSeats = config.NumSeats

	c.Consumer.ProcessEvent(&event)

	candidates := c.State.Pool.Candidates()
	for _, candidate := range candidates {
		c.State.Pool.SetWeight(candidate.Id, 1*c.State.Scale)
	}
}

func (c *commander) ExcludeWithdrawnCandidates(ids []string) {
	excluded := state.MeekCandidates{}

	for _, id := range ids {
		candidate := c.State.Pool.Exclude(id)

		excluded = append(excluded, candidate)
	}

	c.Consumer.ProcessEvent(&CandidatesExcluded{excluded})
}

func (c *commander) PerformPreliminaryElection() {
	numCandidates := c.State.Pool.Count()
	numExcluded := c.State.Pool.ExcludedCount()

	if numCandidates <= (c.State.NumSeats + numExcluded) {
		c.State.Pool.ElectHopeful()
		c.State.ElectedAll = true
		c.Consumer.ProcessEvent(&AllHopefulCandidatesElected{})
	}
}

func (c *commander) IncrementRound() {
	c.State.Round = c.State.Round + 1
	c.State.MeekRound = state.MeekRound{}

	c.Consumer.ProcessEvent(&RoundStarted{c.State.Round})
}

func (c *commander) ExcludeLowestCandidate() {
	lowestCandidates := c.State.Pool.Lowest()

	toExclude := lowestCandidates[0]

	randomUsed := false

	if len(lowestCandidates) > 1 {
		seed := rand.NewSource(time.Now().Unix())
		r := rand.New(seed)
		i := r.Intn(len(lowestCandidates))
		toExclude = lowestCandidates[i]

		randomUsed = true
	}

	c.State.Pool.Exclude(toExclude.Id)
	c.Consumer.ProcessEvent(&LowestCandidateExcluded{lowestCandidates, toExclude, randomUsed})
}

func (c *commander) ExcludeRemainingCandidates() {
	candidates := c.State.Pool.Candidates()

	for _, candidate := range candidates {
		if candidate.Status != state.Elected {
			c.State.Pool.Exclude(candidate.Id)
		}
	}

	c.Consumer.ProcessEvent(&RemainingCandidatesExcluded{})
}

func (c *commander) DistributeVotes() {
	fmt.Println("distributing votes")
	for i := 0; i < c.State.MaxIterations; i++ {
		c.State.MeekRound.Excess = 0

		for _, ballot := range c.State.Ballots {
			value := int64(ballot.Count) * c.State.Scale

			ended := false

			iter := ballot.Ballot.List.Front()
			fmt.Println("c", iter.Value.(string))

			for {
				candidate := c.State.Pool.Candidate(iter.Value.(string))

				if !ended && candidate.Weight > 0 {
					ended = candidate.Status == state.Hopeful

					if ended {
						votes := candidate.Votes + value
						c.State.Pool.SetVotes(candidate.Id, votes)
						value = 0
					} else {
						votes := candidate.Votes + value*candidate.Weight
						c.State.Pool.SetVotes(candidate.Id, votes)
						value = value * (c.State.Scale - candidate.Weight) / c.State.Scale
						fmt.Println(value)
					}
				}

				if iter.Next() == nil {
					break
				}

				iter = iter.Next()
			}

			c.State.MeekRound.Excess = c.State.MeekRound.Excess + value

		}

		c.State.Quota = 1600000

		converged := true
		candidates := c.State.Pool.Candidates()
		for _, candidate := range candidates {
			if candidate.Status == state.Elected {
				temp := c.State.Quota * candidate.Weight

				a := temp / candidate.Votes

				remaineder := temp % candidate.Votes

				if remaineder > 0 {
					a = a + 1
				}

				if a > 1000010 || a < 999990 {
					converged = false
				}

				c.State.Pool.SetWeight(candidate.Id, a)
				fmt.Println(candidate.Name, " weight set to ", a)
			}
		}

		if converged {
			break
		}

	}

	count := 0
	candidates := c.State.Pool.Candidates()
	for _, candidate := range candidates {
		if candidate.Status == state.Hopeful && candidate.Votes > c.State.Quota {
			count = count + 1
			c.State.Pool.Almost(candidate.Id)
		}
	}

	candidates = c.State.Pool.Candidates()
	for _, candidate := range candidates {
		if candidate.Status == state.Almost {
			c.State.Pool.Elect(candidate.Id)
			c.State.MeekRound.AnyElected = true
		}
	}

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
	if c.Error != nil {
		return true
	}

	if !c.State.MeekRound.AnyElected {
		return true
	}

	numElected := c.State.Pool.ElectedCount()

	if numElected >= c.State.NumSeats {
		return true
	}

	return false
}

func (c *commander) Status() CommanderStatus {
	status := CommanderStatus{}
	status.Error = c.Error
	status.Events = c.Consumer.Events()
	status.Elected = c.State.Pool.Elected()

	return status
}
