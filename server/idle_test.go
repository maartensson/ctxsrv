package server

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestIdleTracker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	timeout := time.Second

	track := idleTracker(ctx, cancel, timeout)

	track(nil, http.StateNew)

	track(nil, http.StateClosed)
	track(nil, http.StateNew)
	track(nil, http.StateNew)
	track(nil, http.StateNew)
	track(nil, http.StateClosed)
	track(nil, http.StateClosed)
	track(nil, http.StateClosed)
	track(nil, http.StateClosed)
	track(nil, http.StateClosed)
	track(nil, http.StateClosed)

	select {
	case <-ctx.Done():
	case <-time.After(timeout*2 + time.Millisecond*50):
		t.Fatalf("tracker did not stop")
	}
}
