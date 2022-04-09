package factory

import (
	"errors"
	"strings"

	"github.com/shawntoffel/election"
	"github.com/shawntoffel/meekstv"
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
		return meekstv.NewMeekStv(), nil
	}

	return nil, errors.New("unsupported counting method: " + name)
}

func NewCounter(method string) (election.Counter, error) {
	counterFactory := NewCounterFactory()
	return counterFactory.GetCounter(method)
}
