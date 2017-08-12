package event

import (
	"bytes"
	"fmt"
	"github.com/shawntoffel/electioncounter/counters/stv/meek/state"
	"github.com/shawntoffel/electioncounter/election"
	"github.com/shawntoffel/math"
	"math/rand"
	"time"
)

func (e *Create) Process(state *state.MeekState) (string, error) {
	state.NumSeats = e.NumSeats
	state.Precision = e.Precision
	state.Pool.AddNewCandidates(e.Candidates)
	state.Ballots = e.Ballots.Rollup()

	state.Scale = math.Pow64(int64(10), int64(state.Precision))
	state.MaxIterations = 1000

	buffer := bytes.Buffer{}

	buffer.WriteString("A new Meek STV count has been created")
	buffer.WriteString(fmt.Sprintf("\nCandidates: %d", len(e.Candidates)))
	buffer.WriteString(fmt.Sprintf("\nBallots: %d", len(e.Ballots)))
	buffer.WriteString(fmt.Sprintf("\nSeats: %d", state.NumSeats))
	buffer.WriteString(fmt.Sprintf("\nPrecision: %d", state.Precision))
	buffer.WriteString(fmt.Sprintf("\nMax iterations: %d", state.MaxIterations))

	return buffer.String(), nil
}

func (e *ElectAll) Process(s *state.MeekState) (string, error) {

	candidates := s.Pool.Candidates()

	for _, c := range candidates {
		if c.Status == state.Hopeful {
			s.Pool.Elect(c.Id)
		}
	}

	s.ElectedAll = true

	return "All hopeful candidates have been elected.", nil
}

func (e *ExcludeLowest) Process(s *state.MeekState) (string, error) {

	lowestCandidates := s.Pool.Lowest()

	lowest := lowestCandidates[0]

	buffer := bytes.Buffer{}

	if len(lowestCandidates) > 1 {
		seed := rand.NewSource(time.Now().Unix())
		r := rand.New(seed)
		i := r.Intn(len(lowestCandidates))
		lowest = lowestCandidates[i]

		buffer.WriteString(fmt.Sprintf("Candidates %v are tied for the lowest number of votes.", lowestCandidates))
		buffer.WriteString(fmt.Sprintf("\nCandidate %s was randomly selected to break the tie.", lowest.Name))
	}

	s.Pool.Exclude(lowest.Id)

	buffer.WriteString("\nCandidate '" + lowest.Name + "' has the lowest votes and is excluded.")

	return buffer.String(), nil
}

func (e *ExcludeRemaining) Process(s *state.MeekState) (string, error) {
	candidates := s.Pool.Candidates()

	for _, candidate := range candidates {
		if candidate.Status != state.Elected {
			s.Pool.Exclude(candidate.Id)
		}
	}

	return "All remaining candidates have been excluded.", nil
}

func (e *IncrementRound) Process(s *state.MeekState) (string, error) {
	s.Round = s.Round + 1

	return fmt.Sprintf("Round %d has started.", s.Round), nil
}

func (e *WithdrawlCandidates) Process(state *state.MeekState) (string, error) {
	names := []string{}

	for _, id := range e.Ids {
		c := state.Pool.Candidate(id)
		names = append(names, c.Name)
		state.Pool.Exclude(id)
	}

	return fmt.Sprintf("The following candidates have been excluded: %v", names), nil
}
