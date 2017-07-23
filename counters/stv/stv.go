package stv

import (
	"github.com/shawntoffel/electioncounter/counters"
)

type StvCounter interface {
	counters.Counter
}

type stvCounter struct {
	Quota           int64
	Round           int
	NumberToElect   int
	Ballots         []*list.List
	Changes         counters.Events
	ExpectedVersion int
}
