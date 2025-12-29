package server

import (
	"context"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

func idleTracker(
	ctx context.Context,
	cancel context.CancelFunc,
	timeout time.Duration,
) func(net.Conn, http.ConnState) {
	var activeConns atomic.Int64
	var lastActivity atomic.Int64

	lastActivity.Store(time.Now().Unix())

	go func() {
		ticker := time.NewTicker(max(timeout/10, time.Second))
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				last := time.Unix(lastActivity.Load(), 0)
				if time.Since(last) >= timeout {
					if activeConns.Load() == 0 {
						cancel()
						return
					}
					activeConns.Store(0)
				}
			}
		}
	}()

	return func(_ net.Conn, s http.ConnState) {
		switch s {
		case http.StateNew:
			activeConns.Add(1)
		case http.StateClosed, http.StateHijacked:
			lastActivity.Store(time.Now().Unix())
			activeConns.Add(-1)
		}
	}
}
