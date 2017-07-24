package meek

type MemoryStore map[string]MeekCandidate

type poolStorage struct {
	Store MemoryStore
}

func NewMemoryStorage() PoolStorage {
	memStore := make(MemoryStore)
	return &poolStorage{memStore}
}

func (s *poolStorage) Candidate(id string) MeekCandidate {
	return s.Store[id]
}

func (s *poolStorage) Candidates() []MeekCandidate {
	candidates := []MeekCandidate{}

	for _, c := range s.Store {
		candidates = append(candidates, c)
	}

	return candidates
}

func (s *poolStorage) SaveCandidate(candidate MeekCandidate) {
	s.Store[candidate.Id] = candidate
}

func (s *poolStorage) Exclude(id string) {
	candidate := s.Candidate(id)

	candidate.Weight = 0

	candidate.Status = Excluded

	s.SaveCandidate(candidate)
}
