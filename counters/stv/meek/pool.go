package meek

import (
	"github.com/shawntoffel/electioncounter/counters"
	//"sort"
)

type PoolStorage interface {
	Candidate(id string) MeekCandidate
	Candidates() []MeekCandidate
	SaveCandidate(candidate MeekCandidate)
}

type Pool interface {
	Candidate(id string) MeekCandidate
	Candidates() []MeekCandidate
	AddNewCandidates(candidates counters.Candidates)
	//SetKeepValue(id string, value int64)
	//SetVotes(id string, value int64)
	//SetStatus(id string, status CandidateStatus)
	//SortedCandidates() []Candidate
	//Elected() []Candidate
	//TotalFirstRankCount() int
	//SetFirstRankCount(id string, count int)
}

type pool struct {
	Storage PoolStorage
}

func NewPool(storage PoolStorage) Pool {
	return &pool{storage}
}

func (p *pool) Candidate(id string) MeekCandidate {
	return p.Storage.Candidate(id)
}

func (p *pool) Candidates() []MeekCandidate {
	return p.Storage.Candidates()
}

func (p *pool) AddNewCandidates(candidates counters.Candidates) {
	for _, c := range candidates {
		meekCandidate := MeekCandidate{}
		meekCandidate.Id = c.Id
		meekCandidate.Name = c.Name
		meekCandidate.Weight = 1
		meekCandidate.Status = Hopeful
		meekCandidate.Votes = 0

		p.Storage.SaveCandidate(meekCandidate)
	}
}

/*

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
*/
