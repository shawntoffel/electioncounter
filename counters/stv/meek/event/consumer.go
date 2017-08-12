package event

import (
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
)

type Consumer interface {
	ProcessEvent(event MeekEvent)
}

type consumer struct {
	State state.MeekState
}

func NewConsumer() Consumer {
	c := consumer{}
	c.State = state.NewMeekState{}
	return &c
}

func (c *consumer) ProcessEvent(event MeekEvent) {
	if c.State.Error != nil {
		return
	}

	description, err := event.Process(&c.State)

	if err != nil {
		c.State.Error = err
		return
	}

	counterEvent := election.Event{}
	counterEvent.Description = description

	c.State.Events = append(c.State.Events, counterEvent)
}
