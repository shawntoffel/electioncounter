package counters

import (
	"container/list"
)

func Equal(left *list.List, right *list.List) bool {
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

func Contains(ballots Ballots, ballot *list.List) bool {
	for _, b := range ballots {
		if Equal(b, ballot) {
			return true
		}
	}

	return false
}

func ContainsValue(ballots map[int]Ballot, ballot Ballot) (bool, int) {
	for i, b := range ballots {
		if Equal(b, ballot) {
			return true, i
		}
	}

	return false, 0
}

type RolledUpBallots []RolledUpBallot
type RolledUpBallot struct {
	Count  int
	Ballot Ballot
}

func Rollup(ballots Ballots) []RolledUpBallot {
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

	results := []RolledUpBallot{}

	for i, ballot := range rolledUp {
		result := RolledUpBallot{}

		result.Count = counter[i]
		result.Ballot = ballot

		results = append(results, result)
	}

	return results
}
