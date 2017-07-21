package counters

import (
	"container/list"
)

type MeekStvCounter interface {
	HandleEvent(event CounterEvent)
	Events() []interface{}

	Initialize(config StvConfig)

	SetInitialQuota()
	InitializeVotes()
	UpdateRound()
	HasEnded() bool
	Results() ([]Candidate, Events)
}

type counter struct {
	aggregate

	Quota         int64
	Precision     int
	Round         int
	NumberToElect int
	Ballots       []*list.List
	Pool          Pool
	Scaler        int64
}

func NewMeekStvCounter(events []CounterEvent) MeekStvCounter {
	c := counter{}
	for _, event := range events {
		c.HandleEvent(event)
		c.ExpectedVersion++
	}

	return &c
}

func (state *counter) HandleEvent(event CounterEvent) {
	state.Changes = append(state.Changes, event)

	event.Transition(state)
}

func (state *counter) Events() []interface{} {
	return state.Changes
}

func (state *counter) SetInitialQuota() {
	var numBallots = int64(state.Pool.TotalFirstRankCount())
	var droop = ((numBallots*state.Scaler)/(int64(state.NumberToElect)+1))/state.Scaler*state.Scaler + state.Scaler

	event := QuotaUpdated{}
	event.NewQuota = droop

	state.HandleEvent(&event)
}

func (state *counter) Initialize(config StvConfig) {
	countCreated := CountCreated{}
	countCreated.Candidates = config.Candidates
	countCreated.Precision = config.Precision
	countCreated.NumberToElect = config.NumberToElect

	state.HandleEvent(&countCreated)
	state.RollUpBallots(config.Ballots)
}

func (state *counter) InitializeVotes() {
	candidates := state.Pool.Candidates()

	for _, c := range candidates {
		var votes = int64(c.FirstRankCount) * state.Scaler

		state.UpdateCandidateVotes(c, votes)
		state.Pool.SetKeepValue(c.Id, state.Scaler)

		if state.Pool.Candidate(c.Id).Votes > state.Quota {
			state.ElectCandidate(c)
		}
	}
}

func (state *counter) UpdateCandidateForRound(candidate Candidate) {
	if candidate.Votes > state.Quota {
		var num = state.Quota * candidate.KeepValue
		var newKeepValue = num / candidate.Votes

		var remainder = num % candidate.Votes

		if remainder > 0 {
			newKeepValue += 1
		}

		state.UpdateCandidateKeepValue(candidate, newKeepValue)

		state.ElectCandidate(candidate)
	}

	var newVotes = int64(candidate.FirstRankCount) * state.Pool.Candidate(candidate.Id).KeepValue

	state.UpdateCandidateVotes(candidate, newVotes)
}

func (state *counter) UpdateCandidateVotes(candidate Candidate, votes int64) {
	event := CandidateVotesUpdated{}
	event.Id = candidate.Id
	event.NewVotes = votes

	state.HandleEvent(&event)

}

func (state *counter) UpdateCandidateKeepValue(candidate Candidate, keepValue int64) {
	keepValueUpdated := CandidateKeepValueUpdated{}
	keepValueUpdated.Id = candidate.Id
	keepValueUpdated.NewKeepValue = keepValue

	state.HandleEvent(&keepValueUpdated)
}

func (state *counter) ElectCandidate(candidate Candidate) {
	elected := ElectCandidate{}
	elected.Id = candidate.Id
	state.HandleEvent(&elected)
}

func (state *counter) UpdateRound() {
	state.HandleEvent(&IncrementRound{})

	candidates := state.Pool.Candidates()

	for _, c := range candidates {
		state.UpdateCandidateForRound(c)
	}

	for _, ballot := range state.Ballots {
		if ballot.Len() < 2 {
			continue
		}

		var iter = ballot.Front()

		var topNode = state.Pool.Candidate(iter.Value.(string))

		var multiplier = state.Scaler

		state.DistributeVotes(topNode, iter, multiplier)

	}

	for _, c := range state.Pool.Candidates() {
		if c.Votes > state.Quota {
			state.ElectCandidate(c)
		}
	}

}

func (state *counter) DistributeVotes(firstCandidate Candidate, iter *list.Element, multiplier int64) {

	for {
		if iter.Next() == nil {
			break
		}

		iter = iter.Next()

		var currentCandidate = state.Pool.Candidate(iter.Value.(string))
		var previousCandidate = state.Pool.Candidate(iter.Prev().Value.(string))

		multiplier = multiplier * (state.Scaler - previousCandidate.KeepValue) / state.Scaler

		var votes = currentCandidate.Votes + (multiplier*currentCandidate.KeepValue*int64(firstCandidate.FirstRankCount))/state.Scaler

		state.UpdateCandidateVotes(currentCandidate, votes)

		if iter.Next() == nil {
			var excess = multiplier * (state.Scaler - currentCandidate.KeepValue)

			state.UpdateQuota(excess)
		}
	}

}

func (state *counter) UpdateQuota(excess int64) {

	var numVotes = int64(state.Pool.TotalFirstRankCount())

	var quota = ((numVotes - excess) * state.Scaler / (int64(state.NumberToElect) + 1)) / state.Scaler * state.Scaler

	event := QuotaUpdated{}
	event.NewQuota = quota

	state.HandleEvent(&event)
}

func (state *counter) HasEnded() bool {
	elected := state.Pool.Elected()

	if len(elected) >= state.NumberToElect {
		return true
	}

	for _, c := range elected {
		var frac = float64(state.Quota) / float64(c.Votes)

		if frac < 0.99999 || frac > 1.00001 {
			return false
		}
	}

	return true

}

func (state *counter) Results() ([]Candidate, Events) {
	return state.Pool.Elected(), state.Changes
}

func (state *counter) RollUpBallots(ballots Ballots) {
	for _, ballot := range ballots {
		var first = ballot.Front().Value.(string)

		var candidate = state.Pool.Candidate(first)
		state.Pool.SetFirstRankCount(candidate.Id, candidate.FirstRankCount+1)

		if !Contains(state.Ballots, ballot) {
			state.Ballots = append(state.Ballots, ballot)
		}
	}
}

func Equal(left *list.List, right *list.List) bool {
	if left.Len() != right.Len() {
		return false
	}

	var leftIter = left.Front()
	var rightIter = right.Front()

	var len = left.Len()

	for i := 0; i < len; i++ {
		if leftIter.Value != rightIter.Value {
			return false
		}

		leftIter = leftIter.Next()
		rightIter = rightIter.Next()
	}

	return true
}

func Contains(ballots Ballots, ballot *list.List) bool {
	for _, b := range ballots {
		if Equal(b, ballot) {
			return true
		}
	}

	return false
}
