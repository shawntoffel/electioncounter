package main

type Pool interface {
	SetKeepValue(id string, value int64)
	SetVotes(id string, value int64)
	SetStatus(id string, status CandidateStatus)
	Candidate(id string) Candidate
	Candidates() []Candidate
	Elected() []Candidate
	TotalFirstRankCount() int
}

type pool struct {
	CandidatePool map[string]Candidate
}

func NewPool(candidates []Candidate) Pool {
	pool := pool{}

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
}

func (p *pool) SetVotes(id string, value int64) {
	c := p.Candidate(id)

	c.Votes = value
}

func (p *pool) SetStatus(id string, status CandidateStatus) {
	c := p.Candidate(id)

	c.Status = status
}

func (p *pool) Candidates() []Candidate {
	candidates := []Candidate{}

	for _, c := range p.CandidatePool {
		candidates = append(candidates, c)
	}

	return candidates
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
