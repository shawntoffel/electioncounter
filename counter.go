package main

import (
	"container/list"
)

type Ballot *list.List

type Counter interface {
	HandleEvent(event CounterEvent)
}

type counter struct {
	aggregate

	Quota         int64
	Precision     int
	Round         int
	NumberToElect int
	Ballots       []Ballot
	Pool          Pool
}

func NewCounter(events []CounterEvent) Counter {
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
