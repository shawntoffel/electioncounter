package meek

import (
	"github.com/shawntoffel/electioncounter/counters/stv"
	"github.com/shawntoffel/electioncounter/election"
)

type Meek interface {
	Create(config election.Config)
	WithdrawlCandidates(ids []string)
	HasEnded() bool
	HandleEvent(event MeekEvent)
	Changes() (election.Events, error)
}

type meekStv struct {
	stv.Stv
	Pool      Pool
	Precision int
	Scale     int64
}

func NewMeek(events []MeekEvent) Meek {
	m := meekStv{}
	m.Pool = NewPool()

	for _, event := range events {
		m.HandleEvent(event)
		m.ExpectedVersion++
	}

	return &m
}

func (m *meekStv) Create(config election.Config) {

	if m.Error != nil {
		return
	}

	event := CreateCount{}
	event.Candidates = config.Candidates
	event.Ballots = config.Ballots
	event.Precision = config.Precision
	event.NumberToElect = config.NumberToElect

	m.HandleEvent(&event)
}

func (m *meekStv) WithdrawlCandidates(ids []string) {
	if m.Error != nil {
		return
	}

	event := WithdrawlCandidates{}
	event.Ids = ids

	m.HandleEvent(&event)
}

func (m *meekStv) HasEnded() bool {
	if m.Error != nil {
		return true
	}

	return true
}

func (m *meekStv) Changes() (election.Events, error) {
	return m.Events, m.Error
}

func (m *meekStv) HandleEvent(event MeekEvent) {
	description := event.Transition(m)

	counterEvent := election.Event{}
	counterEvent.Description = description

	m.Events = append(m.Events, counterEvent)
}
