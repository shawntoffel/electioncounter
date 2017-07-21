package main

import (
	"container/list"
)

type Stv interface {
	//	Count(config StvConfig) ([]Candidate, []interface{})
}

type stv struct {
	Counter Counter
}

type StvConfig struct {
	NumberToElect int
	Ballots       []*list.List
	Candidates    []Candidate
	Precision     int
}

func GetScaler(precision int) int64 {
	var scaler = int64(10)
	for i := 0; i < precision; i++ {
		scaler *= 10
	}

	return scaler
}
