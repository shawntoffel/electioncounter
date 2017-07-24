package electioncounter

import (
	"errors"
	"github.com/shawntoffel/electioncounter/counters"
	"github.com/shawntoffel/electioncounter/counters/stv/meek"
	"strings"
)

type CounterFactory interface {
	GetCounter(name string) (counters.Counter, error)
}

type counterFactory struct{}

func NewCounterFactory() CounterFactory {
	return &counterFactory{}
}

func (c *counterFactory) GetCounter(name string) (counters.Counter, error) {

	if strings.EqualFold(name, "meek stv") {
		counter := meek.NewMeekStvCounter(nil)

		return counter, nil

	}

	return nil, errors.New("Unsupported counting method: " + name)
}
