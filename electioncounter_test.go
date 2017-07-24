package electioncounter

import (
	"container/list"
	"fmt"
	"github.com/shawntoffel/electioncounter/counters"
	"testing"
)

func TestMeekStv(t *testing.T) {

	var config = Config{}
	config.Method = "meek stv"

	names := []string{"Alice", "Bob", "Chris", "Don", "Eric"}

	for _, name := range names {
		c := counters.Candidate{}
		c.Id = name
		c.Name = name

		config.Candidates = append(config.Candidates, c)
	}

	var ballots counters.Ballots

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

	var cm = NewElectionCounter()

	var result, err = cm.Count(config)

	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println("Events:", len(result.Events))

	for _, e := range result.Events {
		fmt.Println(e.Description)
	}

	count := len(result.Candidates)
	expectedCount := 3

	if count != expectedCount {
		t.Errorf("Incorrect number of elected candidates. Expected: %d, Got: %d", expectedCount, count)
	}
}
