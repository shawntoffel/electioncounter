package counters

import (
	"container/list"
)

type MeekStvCounter interface {
	StvCounter
}

type counter struct {
	stvCounter

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

func (state *counter) CreateCount(config StvConfig) {
	countCreated := CountCreated{}
	countCreated.Candidates = config.Candidates
	countCreated.Precision = config.Precision
	countCreated.NumberToElect = config.NumberToElect

	state.HandleEvent(&countCreated)
	state.RollUpBallots(config.Ballots)

	state.SetInitialQuota()
	state.InitializeVotes()
}

func (state *counter) InitializeVotes() {
	candidates := state.Pool.Candidates()

	for _, c := range candidates {
		var votes = int64(c.FirstRankCount) * state.Scaler

		state.UpdateCandidateVotes(c, votes)
		state.Pool.SetKeepValue(c.Id, state.Scaler)
	}
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

func (state *counter) Events() Events {
	return state.Changes
}

func (state *counter) Results() ([]Candidate, Events) {
	return state.Pool.Elected(), state.Changes
}

func (state *counter) SetInitialQuota() {
	initialVoteCount := int64(state.Pool.TotalFirstRankCount())
	numberToElect := int64(state.NumberToElect)

	var droopQuota = ((initialVoteCount*state.Scaler)/(numberToElect+1))/state.Scaler*state.Scaler + state.Scaler

	event := QuotaUpdated{}
	event.NewQuota = droopQuota

	state.HandleEvent(&event)
}

func (state *counter) UpdateCandidateForRound(candidate Candidate) {
	if candidate.Votes > state.Quota {
		num := state.Quota * candidate.KeepValue
		newKeepValue := num / candidate.Votes

		remainder := num % candidate.Votes

		if remainder > 0 {
			newKeepValue += 1
		}

		state.UpdateCandidateKeepValue(candidate, newKeepValue)
	}

	newVotes := int64(candidate.FirstRankCount) * state.Pool.Candidate(candidate.Id).KeepValue

	state.UpdateCandidateVotes(candidate, newVotes)
}

func (state *counter) UpdateCandidateVotes(candidate Candidate, votes int64) {
	if candidate.Votes == votes {
		return
	}

	event := CandidateVotesUpdated{}
	event.Id = candidate.Id
	event.NewVotes = votes

	state.HandleEvent(&event)

	state.ElectCandidateAboveQuota(candidate)
}

func (state *counter) GiveVotesToCandidate(candidate Candidate, votes int64) {
	event := CandidateVotesReceived{}
	event.Id = candidate.Id
	event.ReceivedVotes = votes

	state.HandleEvent(&event)

	state.ElectCandidateAboveQuota(candidate)
}

func (state *counter) ElectCandidateAboveQuota(candidate Candidate) {
	updatedCandidate := state.Pool.Candidate(candidate.Id)
	isElected := updatedCandidate.Status == Elected

	if !isElected && updatedCandidate.Votes > state.Quota {
		state.ElectCandidate(candidate)
	}
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

		iter := ballot.Front()

		firstPriorityCandidate := state.Pool.Candidate(iter.Value.(string))

		state.DistributeVotes(firstPriorityCandidate, iter)
	}
}

func (state *counter) DistributeVotes(firstCandidate Candidate, iter *list.Element) {

	state.HandleEvent(&DistributeVotes{})
	proportion := state.Scaler
	initialVoteCount := int64(firstCandidate.FirstRankCount)

	for {
		if iter.Next() == nil {
			break
		}

		iter = iter.Next()

		currentCandidate := state.Pool.Candidate(iter.Value.(string))
		previousCandidate := state.Pool.Candidate(iter.Prev().Value.(string))

		proportion = proportion * (state.Scaler - previousCandidate.KeepValue) / state.Scaler

		votesToKeep := (proportion * currentCandidate.KeepValue * initialVoteCount) / state.Scaler

		state.GiveVotesToCandidate(currentCandidate, votesToKeep)

		if iter.Next() == nil {
			var excess = proportion * (state.Scaler - currentCandidate.KeepValue)

			state.UpdateQuota(excess)
		}
	}

}

func (state *counter) UpdateQuota(excess int64) {

	initialVoteCount := int64(state.Pool.TotalFirstRankCount())
	numberToElect := int64(state.NumberToElect)

	quota := ((initialVoteCount - excess) * state.Scaler / (numberToElect + 1)) / state.Scaler * state.Scaler

	if state.Quota == quota {
		return
	}

	event := QuotaUpdated{}
	event.NewQuota = quota

	state.HandleEvent(&event)
}

func (state *counter) RollUpBallots(ballots Ballots) {
	for _, ballot := range ballots {
		first := ballot.Front().Value.(string)

		candidate := state.Pool.Candidate(first)
		state.Pool.SetFirstRankCount(candidate.Id, candidate.FirstRankCount+1)

		if !Contains(state.Ballots, ballot) {
			state.Ballots = append(state.Ballots, ballot)
		}
	}
}
