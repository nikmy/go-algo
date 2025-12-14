package netx

import (
	"errors"
	"testing"

	"github.com/nikmy/algo/container/deque"
	"github.com/nikmy/algo/testx/faulty"
)

func newFaultyChannel(t *testing.T, cfg NetworkConfig) faultyChannel {
	c := faultyChannel{
		sent:    new(deque.Deque[AssertableRemoteCall]),
		reorder: faulty.NewController(t, 1),
		srcLoss: faulty.NewController(t, 2),
		dup:     faulty.NewController(t, 4),
	}

	c.reorder.SetFaultProbability(cfg.ReorderingProb)
	c.srcLoss.SetFaultProbability(cfg.SendLossProb)
	c.dstLoss.SetFaultProbability(cfg.RecvLossProb)
	c.dup.SetFaultProbability(cfg.DuplicateProb)

	return c
}

type faultyChannel struct {
	sent    *deque.Deque[AssertableRemoteCall]
	reorder *faulty.Controller
	srcLoss *faulty.Controller
	dstLoss *faulty.Controller
	dup     *faulty.Controller
	world   *faulty.Controller
}

func (c faultyChannel) looseAll() {
	c.sent.Clear()
}

func (c faultyChannel) send(msg AssertableRemoteCall) error {
	if c.srcLoss.Fault() {
		return errors.New("connection failure")
	}

	c.sent.PushBack(msg)
	if c.dup.Fault() {
		c.sent.PushBack(msg)
	}

	return nil
}

func (c faultyChannel) pull() *AssertableRemoteCall {
	if c.sent.IsEmpty() {
		return nil
	}

	next := c.sent.PopFront()
	if c.dstLoss.Fault() {
		return nil
	}

	if c.sent.Len() > 1 && c.reorder.Fault() {
		sec := c.sent.PopFront()
		c.sent.PushFront(next)
		next = sec
	}

	return &next
}

