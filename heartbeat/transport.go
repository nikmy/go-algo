package heartbeat

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/nikmy/algo/fp"
)

type transport[Msg any] interface {
	Dial(ctx context.Context, msg *Msg) error
}

type wrapper[Msg any] struct {
	next atomic.Pointer[Msg]
	impl transport[Msg]
}

func (w *wrapper[Msg]) Dial(_ context.Context, msg *Msg) error {
	w.next.Store(msg)
	return nil
}

func Start[Msg any](
	ctx context.Context,
	interval time.Duration,
	baseCtx func(context.Context) context.Context,
	baseTransport transport[Msg],
	errHandler func(error),
) *wrapper[Msg] {
	w := &wrapper[Msg]{impl: baseTransport}
	if baseCtx == nil {
		baseCtx = fp.Id[context.Context]
	}
	go func() {
		tick := time.NewTicker(interval)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				err := w.impl.Dial(baseCtx(ctx), w.next.Swap(nil))
				if err != nil {
					errHandler(err)
				}
			}
		}
	}()
	return w
}
