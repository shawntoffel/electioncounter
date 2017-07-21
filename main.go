package main

import (
	"container/list"
	"fmt"
)

func main() {

	var config = StvConfig{}

	names := []string{"Alice", "Bob", "Chris", "Don", "Eric"}

	for _, name := range names {
		c := Candidate{}
		c.Id = name
		c.Status = Hopeful

		config.Candidates = append(config.Candidates, c)
	}

	var ballots []*list.List

	for i := 0; i < 28; i++ {
		var ballot = list.New()
		ballot.PushBack("Alice")
		ballot.PushBack("Bob")
		ballot.PushBack("Chris")
		ballots = append(ballots, ballot)
	}

	for i := 0; i < 26; i++ {
		var ballot = list.New()
		ballot.PushBack("Bob")
		ballot.PushBack("Alice")
		ballot.PushBack("Chris")
		ballots = append(ballots, ballot)
	}

	for i := 0; i < 3; i++ {
		var ballot = list.New()
		ballot.PushBack("Chris")
		ballots = append(ballots, ballot)
	}

	for i := 0; i < 2; i++ {
		var ballot = list.New()
		ballot.PushBack("Don")
		ballots = append(ballots, ballot)
	}

	var ballot = list.New()
	ballot.PushBack("Eric")
	ballots = append(ballots, ballot)

	config.Ballots = ballots

	config.NumberToElect = 3
	config.Precision = 6

	var cm = NewMeekStv()

	cm.Initialize(config)

	var candidates, events = cm.Run()

	for _, c := range candidates {
		fmt.Println(c.Id)
	}

	fmt.Print(events)
}
