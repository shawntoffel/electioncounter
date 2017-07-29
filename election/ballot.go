package election

import (
	"container/list"
)

type Ballots []Ballot
type Ballot struct {
	*list.List
}

type RolledUpBallots []RolledUpBallot
type RolledUpBallot struct {
	Count  int
	Ballot Ballot
}

func (left Ballot) Equal(right Ballot) bool {
	if left.Len() != right.Len() {
		return false
	}

	leftIter := left.Front()
	rightIter := right.Front()

	for i := 0; i < left.Len(); i++ {
		if leftIter.Value != rightIter.Value {
			return false
		}

		leftIter = leftIter.Next()
		rightIter = rightIter.Next()
	}

	return true
}

func (ballots Ballots) Contains(ballot Ballot) bool {
	for _, b := range ballots {
		if b.Equal(ballot) {
			return true
		}
	}

	return false
}

func ContainsValue(ballots map[int]Ballot, ballot Ballot) (bool, int) {
	for i, b := range ballots {
		if b.Equal(ballot) {
			return true, i
		}
	}

	return false, 0
}

func (ballots Ballots) Rollup() RolledUpBallots {
	counter := make(map[int]int)
	rolledUp := make(map[int]Ballot)

	for _, ballot := range ballots {
		contains, index := ContainsValue(rolledUp, ballot)

		if contains {
			counter[index] = counter[index] + 1
		} else {
			rolledUp[len(rolledUp)+1] = ballot
		}
	}

	results := RolledUpBallots{}

	for i, ballot := range rolledUp {
		result := RolledUpBallot{}

		result.Count = counter[i]
		result.Ballot = ballot

		results = append(results, result)
	}

	return results
}
