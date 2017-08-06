package state

import (
	"github.com/shawntoffel/electioncounter/election"
	"github.com/shawntoffel/memorystorage"
	"sort"
)

type Pool interface {
	Candidate(id string) MeekCandidate
	SetVotes(id string, votes int64)
	Lowest() MeekCandidates
	Candidates() MeekCandidates
	Count() int
	ElectedCount() int
	ExcludedCount() int
	Elect(id string)
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

func (p *pool) SetVotes(id string, votes int64) {
	candidate := p.Candidate(id)

	candidate.Votes = votes

	p.Storage.Set(candidate.Id, candidate)
}

func (p *pool) Candidates() MeekCandidates {
	candidates := MeekCandidates{}

	list := p.Storage.List()
	for _, candidate := range list {
		candidates = append(candidates, candidate.(MeekCandidate))
	}

	return candidates
}

func (p *pool) CandidatesWithStatus(status CandidateStatus) MeekCandidates {
	candidates := p.Candidates()
	result := MeekCandidates{}

	for _, candidate := range candidates {
		if candidate.Status == status {
			result = append(result, candidate)
		}
	}

	return result
}

func (p *pool) Count() int {
	return len(p.Storage.List())
}

func (p *pool) ExcludedCount() int {
	return len(p.CandidatesWithStatus(Excluded))
}

func (p *pool) ElectedCount() int {
	return len(p.CandidatesWithStatus(Elected))
}

func (p *pool) Elect(id string) {
	candidate := p.Candidate(id)

	candidate.Status = Elected

	p.Storage.Set(candidate.Id, candidate)
}

func (p *pool) Lowest() MeekCandidates {
	candidates := p.Candidates()

	sort.Sort(ByVotes(candidates))

	lowest := MeekCandidates{}

	for _, candidate := range candidates {
		if len(lowest) > 0 && candidate.Votes != lowest[0].Votes {
			break
		}

		lowest = append(lowest, candidate)
	}

	return lowest
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
