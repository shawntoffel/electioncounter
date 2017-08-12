package event

import (
	"github.com/shawntoffel/electioncounter/election"
)

type Consumer interface {
	ProcessEvent(event MeekEvent)
}

type consumer struct {
	Events election.Events
}

func NewConsumer() Consumer {
	c := consumer{}
	c.Events = election.Events{}

	return &c
}

func (c *consumer) ProcessEvent(meekEvent MeekEvent) {
	event := meekEvent.Process()

	c.Events = append(c.Events, event)
}
