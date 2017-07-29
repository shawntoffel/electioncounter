package event

import (
	"github.com/shawntoffel/electioncounter/counters/stv/meek/event/events"
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
	"github.com/shawntoffel/electioncounter/election"
)

type MeekEventProcessor interface {
	Create(config election.Config)
	WithdrawlCandidates(ids []string)
	HasEnded() bool
	HandleEvent(event MeekEvent)
	Changes() (election.Events, error)
}

type meekEventProcessor struct {
	Error error
	State *state.MeekState
}

func NewMeekEventProcessor(events []MeekEvent) MeekEventProcessor {
	s := meekEventProcessor{}
	s.State = &state.MeekState{}
	s.State.Pool = state.NewPool()

	for _, event := range events {
		s.HandleEvent(event)
		s.State.ExpectedVersion++
	}

	return &s
}

func (s *meekEventProcessor) Create(config election.Config) {

	if s.Error != nil {
		return
	}

	event := events.Create{}
	event.Candidates = config.Candidates
	event.Ballots = config.Ballots
	event.Precision = config.Precision
	event.NumberToElect = config.NumberToElect

	s.HandleEvent(&event)
}

func (s *meekEventProcessor) WithdrawlCandidates(ids []string) {
	if s.Error != nil {
		return
	}

	event := events.WithdrawlCandidates{}
	event.Ids = ids

	s.HandleEvent(&event)
}

func (s *meekEventProcessor) HasEnded() bool {
	if s.Error != nil {
		return true
	}

	return true
}

func (s *meekEventProcessor) Changes() (election.Events, error) {
	return s.State.Events, s.Error
}
