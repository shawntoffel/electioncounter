package events

import (
	"bytes"
	"fmt"
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
	"math/rand"
	"time"
)

type ExcludeLowest struct{}

func (e *ExcludeLowest) Transition(s *state.MeekState) (string, error) {

	lowestCandidates := s.Pool.Lowest()

	lowest := lowestCandidates[0]

	buffer := bytes.Buffer{}

	if len(lowestCandidates) > 1 {
		seed := rand.NewSource(time.Now().Unix())
		r := rand.New(seed)
		i := r.Intn(len(lowestCandidates))
		lowest = lowestCandidates[i]

		buffer.WriteString(fmt.Sprintf("Candidates %v are tied for the lowest number of votes.", lowestCandidates))
		buffer.WriteString(fmt.Sprintf("Candidate %s was randomly selected to break the tie.", lowest.Name))
	}

	s.Pool.Exclude(lowest.Id)

	return "Candidate '" + lowest.Name + "' has the lowest votes and is excluded.", nil
}
