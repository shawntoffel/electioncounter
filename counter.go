package main

import (
	"container/list"
)

type Ballot *list.List

type Counter struct {
	Aggregate

	NumberOfWinners int
	Ballots         []Ballot
	Candidates      []Candidate
}
