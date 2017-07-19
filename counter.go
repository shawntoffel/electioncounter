package main

import (
	"container/list"
)

type Ballot *list.List

type Counter interface {
	//Aggregate
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

func (state *counter) InitializeCounter(event CountCreated) {
	state.NumberToElect = event.NumberToElect
	state.Ballots = event.Ballots

	state.Precision = event.Precision
	state.Pool = NewPool(event.Candidates)
}

func (state *counter) HandleEvent(event CounterEvent) {
	state.Changes = append(state.Changes, event)

	state.RunTask(event.Task)
}

func (state *counter) RunTask(task EventTask) {
	task(state)
}
