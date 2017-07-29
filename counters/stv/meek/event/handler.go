package event

import (
	"github.com/shawntoffel/electioncounter/election"
)

func (m *meekEventProcessor) HandleEvent(event MeekEvent) {
	description := event.Transition(m.State)

	counterEvent := election.Event{}
	counterEvent.Description = description

	m.State.Events = append(m.State.Events, counterEvent)
}
