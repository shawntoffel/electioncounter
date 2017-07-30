package state

import (
	"github.com/shawntoffel/electioncounter/election"
	"github.com/shawntoffel/memorystorage"
)

type Pool interface {
	Candidate(id string) MeekCandidate
	Candidates() MeekCandidates
	Count() int
	ExcludedCount() int
	ElectCandidate(id string)
	AddNewCandidates(candidates election.Candidates)
	Exclude(id string)
}

type pool struct {
	Storage memorystorage.MemoryStorage
}

func NewPool() Pool {
	p := pool{}
	p.Storage = memorystorage.NewMemoryStorage()
	return &p
}

func (p *pool) Candidate(id string) MeekCandidate {
	return p.Storage.Get(id).(MeekCandidate)
}

func (p *pool) Candidates() MeekCandidates {
	candidates := MeekCandidates{}

	list := p.Storage.List()
	for _, candidate := range list {
		candidates = append(candidates, candidate.(MeekCandidate))
	}

	return candidates
}

func (p *pool) Count() int {
	return len(p.Storage.List())
}

func (p *pool) ExcludedCount() int {
	count := 0
	candidates := p.Candidates()

	for _, candidate := range candidates {
		if candidate.Status == Excluded {
			count++
		}
	}

	return count
}

func (p *pool) ElectCandidate(id string) {
	candidate := p.Candidate(id)

	candidate.Status = Elected

	p.Storage.Set(candidate.Id, candidate)
}

func (p *pool) AddNewCandidates(candidates election.Candidates) {
	for _, c := range candidates {
		meekCandidate := MeekCandidate{}
		meekCandidate.Id = c.Id
		meekCandidate.Name = c.Name
		meekCandidate.Weight = 1
		meekCandidate.Status = Hopeful
		meekCandidate.Votes = 0

		p.Storage.Set(meekCandidate.Id, meekCandidate)
	}
}

func (p *pool) Exclude(id string) {
	candidate := p.Candidate(id)
	candidate.Weight = 0
	candidate.Status = Excluded
	p.Storage.Set(candidate.Id, candidate)
}
