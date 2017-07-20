package main

type Pool interface {
	SetKeepValue(id string, value int64)
	SetVotes(id string, value int64)
	SetStatus(id string, status CandidateStatus)
	Candidate(id string) Candidate
}

type pool struct {
	Candidates map[string]Candidate
}

func NewPool(candidates []Candidate) Pool {
	pool := pool{}

	for _, c := range candidates {
		pool.Candidates[c.Id] = c
	}

	return &pool
}

func (p *pool) Candidate(id string) Candidate {
	return p.Candidates[id]
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
