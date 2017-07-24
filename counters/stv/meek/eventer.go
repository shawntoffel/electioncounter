package meek

import (
	"github.com/shawntoffel/electioncounter/counters"
)

func createCount(state *meekStvCounter, counterConfig counters.CounterConfig) {
	event := CreateCount{}
	event.Candidates = counterConfig.Candidates
	event.Ballots = counterConfig.Ballots
	event.Precision = counterConfig.Precision
	event.NumberToElect = counterConfig.NumberToElect

	HandleEvent(state, &event)
}

func withdrawlCandidates(state *meekStvCounter, ids []string) {
	event := WithdrawlCandidates{}
	event.Ids = ids

	HandleEvent(state, &event)
}

func HandleEvent(state *meekStvCounter, event MeekEvent) {

	description := event.Transition(state)

	counterEvent := counters.Event{}
	counterEvent.Description = description

	state.Changes = append(state.Changes, counterEvent)
}
