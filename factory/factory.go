package factory

import (
	"errors"
	"github.com/shawntoffel/election"
	"github.com/shawntoffel/meekstv"
	"strings"
)

type CounterFactory interface {
	GetCounter(name string) (election.Counter, error)
}

type counterFactory struct{}

func NewCounterFactory() CounterFactory {
	return &counterFactory{}
}

func (c *counterFactory) GetCounter(name string) (election.Counter, error) {

	if strings.EqualFold(name, "meekstv") {
		counter := meekstv.NewMeekStv()

		return counter, nil
	}

	return nil, errors.New("Unsupported counting method: " + name)
}

func NewCounter(method string) (election.Counter, error) {
	counterFactory := NewCounterFactory()

	return counterFactory.GetCounter(method)
}
