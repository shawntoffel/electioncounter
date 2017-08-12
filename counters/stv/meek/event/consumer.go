package event

import (
	"github.com/shawntoffel/electioncounter/election"
)

type Consumer interface {
	ProcessEvent(event MeekEvent)
	Events() election.Events
}

type consumer struct {
	events election.Events
}

func NewConsumer() Consumer {
	c := consumer{}
	c.events = election.Events{}

	return &c
}

func (c *consumer) ProcessEvent(meekEvent MeekEvent) {
	event := meekEvent.Process()

	c.events = append(c.events, event)
}

func (c *consumer) Events() election.Events {
	return c.events
}
