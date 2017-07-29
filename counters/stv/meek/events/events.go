package events

import ()

type MeekEvent interface {
	Transition(m *meekState) string
}
