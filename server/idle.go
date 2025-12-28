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

	lastActivity.Store(time.Now().UnixNano())

	go func() {
		ticker := time.NewTicker(max(timeout/10, time.Second))
		defer ticker.Stop()
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if activeConns.Load() == 0 {
					last := time.Unix(0, lastActivity.Load())
					if time.Since(last) >= timeout {
						return
					}
				}
			}
		}
	}()

	return func(_ net.Conn, s http.ConnState) {
		switch s {
		case http.StateNew:
			activeConns.Add(1)
		case http.StateClosed, http.StateHijacked:
			lastActivity.Store(time.Now().UnixNano())
			activeConns.Add(-1)
		}
	}
}
