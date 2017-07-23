package counters

import (
	"sort"
)

type Pool interface {
	SetKeepValue(id string, value int64)
	SetVotes(id string, value int64)
	SetStatus(id string, status CandidateStatus)
	Candidate(id string) Candidate
	Candidates() []Candidate
	SortedCandidates() []Candidate
	Elected() []Candidate
	TotalFirstRankCount() int
	SetFirstRankCount(id string, count int)
}

type pool struct {
	CandidatePool map[string]Candidate
}

func NewPool(candidates []Candidate) Pool {
	pool := pool{}
	pool.CandidatePool = make(map[string]Candidate)

	for _, c := range candidates {
		pool.CandidatePool[c.Id] = c
	}

	return &pool
}

func (p *pool) Candidate(id string) Candidate {
	return p.CandidatePool[id]
}

func (p *pool) SetKeepValue(id string, value int64) {
	c := p.Candidate(id)

	c.KeepValue = value

	p.CandidatePool[id] = c
}

func (p *pool) SetVotes(id string, value int64) {
	c := p.Candidate(id)

	c.Votes = value

	p.CandidatePool[id] = c
}

func (p *pool) SetStatus(id string, status CandidateStatus) {
	c := p.Candidate(id)

	c.Status = status

	p.CandidatePool[id] = c
}

func (p *pool) Candidates() []Candidate {
	candidates := []Candidate{}

	for _, c := range p.CandidatePool {
		candidates = append(candidates, c)
	}

	return candidates
}

func (p *pool) SortedCandidates() []Candidate {
	candidates := Candidates{}

	for _, c := range p.CandidatePool {
		candidates = append(candidates, c)
	}

	sort.Sort(candidates)

	return candidates
}

type Candidates []Candidate

func (c Candidates) Len() int {
	return len(c)
}
func (c Candidates) Less(i, j int) bool {
	return c[i].Votes < c[j].Votes
}
func (c Candidates) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (p *pool) Elected() []Candidate {
	candidates := []Candidate{}

	for _, c := range p.CandidatePool {
		if c.Status == Elected {
			candidates = append(candidates, c)
		}
	}

	return candidates
}

func (p *pool) TotalFirstRankCount() int {
	count := 0

	for _, c := range p.CandidatePool {
		count += c.FirstRankCount
	}

	return count
}

func (p *pool) SetFirstRankCount(id string, count int) {
	c := p.Candidate(id)

	c.FirstRankCount = count

	p.CandidatePool[id] = c
}
