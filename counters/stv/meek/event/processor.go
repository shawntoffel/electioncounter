package event

import (
	"github.com/shawntoffel/electioncounter/counters/stv/meek/event/events"
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
	"github.com/shawntoffel/electioncounter/election"
)

type MeekEventProcessor interface {
	Create(config election.Config)
	ExcludeWithdrawnCandidates(ids []string)
	IncrementRound()
	PerformPreliminaryElection()
	ExcludeRemainingCandidates()
	ExcludeLowestCandidate()
	HasEnded() bool
	DistributeVotes()
	RoundHasEnded() bool
	Changes() (election.Events, error)
}

type meekEventProcessor struct {
	Error      error
	ElectedAll bool
	State      *state.MeekState
}

func NewMeekEventProcessor(events []MeekEvent) MeekEventProcessor {
	s := meekEventProcessor{}
	s.State = state.NewMeekState()

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

func (s *meekEventProcessor) ExcludeWithdrawnCandidates(ids []string) {
	event := events.WithdrawlCandidates{}
	event.Ids = ids

	s.handleEvent(&event)
}

func (s *meekEventProcessor) PerformPreliminaryElection() {
	numCandidates := s.State.Pool.Count()
	numExcluded := s.State.Pool.ExcludedCount()

	if numCandidates <= (s.State.NumSeats + numExcluded) {
		event := events.ElectAll{}

		s.handleEvent(&event)
	}
}

func (s *meekEventProcessor) HasEnded() bool {
	if s.Error != nil {
		return true
	}

	if s.State.ElectedAll {
		return true
	}

	numElected := s.State.Pool.ElectedCount()

	return numElected == s.State.NumSeats
}

func (s *meekEventProcessor) ExcludeRemainingCandidates() {
	event := events.ExcludeRemaining{}

	s.handleEvent(&event)
}

func (s *meekEventProcessor) ExcludeLowestCandidate() {
	event := events.ExcludeLowest{}

	s.handleEvent(&event)
}

func (s *meekEventProcessor) IncrementRound() {
	s.handleEvent(&events.IncrementRound{})
}

func (s *meekEventProcessor) DistributeVotes() {

}

func (s *meekEventProcessor) RoundHasEnded() bool {

	if !s.State.Round.AnyElected {
		return true
	}

	numElected := s.State.Pool.ElectedCount()

	if numElected >= s.State.NumSeats {
		return true
	}

	return false
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
