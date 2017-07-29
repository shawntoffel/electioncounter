package events

import (
	"github.com/shawntoffel/electioncounter/election"
)

func (state *meekState) HandleEvent(event MeekEvent) {
	description := event.Transition(state)

	counterEvent := election.Event{}
	counterEvent.Description = description

	state.Events = append(state.Events, counterEvent)
}
