package synctest

import "sync"

var (
	F faultInjector = nopFaultInjector{}
	m sync.Mutex
)

type faultInjector interface {
	Fault() bool
}

type nopFaultInjector struct{}

func (nopFaultInjector) Fault() bool {
	return false
}
