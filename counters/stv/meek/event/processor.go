package event

import (
	"github.com/shawntoffel/electioncounter/counters/stv/meek/event/events"
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
	"github.com/shawntoffel/electioncounter/election"
)

type MeekEventProcessor interface {
	Create(config election.Config)
	WithdrawlCandidates(ids []string)
	CountInitialVotes()
	HasEnded() bool
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
		s.handleEvent(event)
		s.State.ExpectedVersion++
	}

	return &s
}

func (s *meekEventProcessor) Create(config election.Config) {
	event := events.Create{}
	event.Candidates = config.Candidates
	event.Ballots = config.Ballots
	event.Precision = config.Precision
	event.NumSeats = config.NumSeats

	s.handleEvent(&event)
}

func (s *meekEventProcessor) WithdrawlCandidates(ids []string) {
	event := events.WithdrawlCandidates{}
	event.Ids = ids

	s.handleEvent(&event)
}

func (s *meekEventProcessor) CountInitialVotes() {
	numCandidates := s.State.Pool.Count()

	if numCandidates < s.State.NumSeats {
		event := events.ElectAll{}

		s.handleEvent(&event)
	}
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

func (m *meekEventProcessor) handleEvent(event MeekEvent) {
	if m.Error != nil {
		return
	}

	description, err := event.Transition(m.State)

	if err != nil {
		m.Error = err
		return
	}

	counterEvent := election.Event{}
	counterEvent.Description = description

	m.State.Events = append(m.State.Events, counterEvent)
}
