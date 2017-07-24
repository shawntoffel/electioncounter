package meek

import (
	"github.com/shawntoffel/electioncounter/counters"
)

func createCount(state *meekStvCounter, counterConfig counters.CounterConfig) {
	createCount := CreateCount{}
	createCount.Candidates = counterConfig.Candidates
	createCount.Ballots = counterConfig.Ballots
	createCount.Precision = counterConfig.Precision
	createCount.NumberToElect = counterConfig.NumberToElect

	HandleEvent(state, &createCount)
}

func withdrawlCandidates(state *meekStvCounter, ids []string) {

}

func HandleEvent(state *meekStvCounter, event MeekEvent) {

	description := event.Transition(state)

	counterEvent := counters.Event{}
	counterEvent.Description = description

	state.Changes = append(state.Changes, counterEvent)
}
