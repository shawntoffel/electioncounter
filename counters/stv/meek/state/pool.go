package state

import (
	"github.com/shawntoffel/electioncounter/election"
	"sort"
)

type Storage map[string]*MeekCandidate

type Pool interface {
	Candidate(id string) *MeekCandidate
	SetVotes(id string, votes int64)
	Lowest() MeekCandidates
	Candidates() MeekCandidates
	Count() int
	Elected() MeekCandidates
	ElectedCount() int
	ExcludedCount() int
	Elect(id string)
	Almost(id string)
	ElectHopeful()
	AddNewCandidates(candidates election.Candidates)
	Exclude(id string) *MeekCandidate
	SetWeight(id string, weight int64)
}

type pool struct {
	Storage Storage
}

func NewPool() Pool {
	p := pool{}
	p.Storage = make(Storage)
	return &p
}

func (p *pool) Candidate(id string) *MeekCandidate {
	return p.Storage[id]
}

func (p *pool) SetVotes(id string, votes int64) {
	candidate := p.Candidate(id)

	candidate.Votes = votes

	p.Storage[candidate.Id] = candidate
}

func (p *pool) SetWeight(id string, weight int64) {
	candidate := p.Candidate(id)

	candidate.Weight = weight

	p.Storage[candidate.Id] = candidate
}

func (p *pool) Candidates() MeekCandidates {
	candidates := MeekCandidates{}

	list := List(p.Storage)
	for _, candidate := range list {
		candidates = append(candidates, candidate)
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
	return len(List(p.Storage))
}

func (p *pool) ExcludedCount() int {
	return len(p.CandidatesWithStatus(Excluded))
}

func (p *pool) Elected() MeekCandidates {
	return p.CandidatesWithStatus(Elected)
}

func (p *pool) ElectedCount() int {
	return len(p.CandidatesWithStatus(Elected))
}

func (p *pool) Elect(id string) {
	candidate := p.Candidate(id)

	candidate.Status = Elected

	p.Storage[candidate.Id] = candidate
}

func (p *pool) Almost(id string) {
	candidate := p.Candidate(id)

	candidate.Status = Almost

	p.Storage[candidate.Id] = candidate
}

func (p *pool) ElectHopeful() {
	candidates := p.Candidates()

	for _, candidate := range candidates {
		if candidate.Status == Hopeful {
			p.Elect(candidate.Id)
		}
	}
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

		p.Storage[c.Id] = &meekCandidate
	}
}

func (p *pool) Exclude(id string) *MeekCandidate {
	candidate := p.Candidate(id)
	candidate.Weight = 0
	candidate.Status = Excluded
	p.Storage[candidate.Id] = candidate

	return p.Candidate(id)
}

func List(storage Storage) MeekCandidates {
	list := MeekCandidates{}

	for _, value := range storage {
		list = append(list, *value)
	}

	return list
}
