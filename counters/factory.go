package counters

import (
	"errors"
	"strings"
)

type CounterFactory interface {
	GetCounter(name string) (Counter, error)
}

type counterFactory struct{}

func NewCounterFactory() CounterFactory {
	return &counterFactory{}
}

func (c *counterFactory) GetCounter(name string) (Counter, error) {

	if strings.EqualFold(name, "meek stv") {

	}

	return nil, errors.New("Unsupported counting method: " + name)
}
