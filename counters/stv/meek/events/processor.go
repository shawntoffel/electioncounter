package events

import (
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

func NewMeekEventProcessor(events []MeekEvent) MeekEventProcessor {
	s := meekState{}
	s.Pool = state.NewPool()

	for _, event := range events {
		s.HandleEvent(event)
		s.ExpectedVersion++
	}

	return &s
}

func (s *meekState) Create(config election.Config) {

	if s.Error != nil {
		return
	}

	event := Create{}
	event.Candidates = config.Candidates
	event.Ballots = config.Ballots
	event.Precision = config.Precision
	event.NumberToElect = config.NumberToElect

	s.HandleEvent(&event)
}

func (s *meekState) WithdrawlCandidates(ids []string) {
	if s.Error != nil {
		return
	}

	event := WithdrawlCandidates{}
	event.Ids = ids

	s.HandleEvent(&event)
}

func (s *meekState) HasEnded() bool {
	if s.Error != nil {
		return true
	}

	return true
}

func (s *meekState) Changes() (election.Events, error) {
	return s.Events, s.Error
}
