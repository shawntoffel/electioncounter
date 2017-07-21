package counters

import (
	"container/list"
)

type Ballots []*list.List

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
