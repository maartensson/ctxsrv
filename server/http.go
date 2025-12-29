package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

func HTTP(
	ctx context.Context,
	log *slog.Logger,
	ln net.Listener,
	handler http.Handler,
	opts ...HTTPOption,
) error {
	if ctx == nil {
		return errors.New("nil context")
	}
	if log == nil {
		return errors.New("nil logger")
	}
	if ln == nil {
		return errors.New("nil listener")
	}
	if handler == nil {
		return errors.New("nil handler")
	}

	cfg := defaultHTTPConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var connStateFn func(net.Conn, http.ConnState)
	if cfg.idleTimeout > 0 {
		connStateFn = idleTracker(ctx, cancel, cfg.idleTimeout)
	}

	srv := &http.Server{
		Handler:           handler,
		BaseContext:       func(ln net.Listener) context.Context { return ctx },
		ConnContext:       cfg.connContext,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      0, // MUST be 0 for streaming
		IdleTimeout:       0,
		ConnState:         connStateFn,
		ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelError),
	}
	defer srv.Close()

	go func() {
		defer cancel()

		log.Info("HTTP server listening", "addr", ln.Addr().String())

		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server failed", "error", err)
		}
	}()

	<-ctx.Done()

	log.Warn("shutting down server")

	shutdownCtx, shutDownCancel := context.WithTimeout(ctx, cfg.shutdownTimeout)
	defer shutDownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Warn("server successfully shut down")

	return nil
}

type HTTPOption func(*httpConfig)

type httpConfig struct {
	idleTimeout     time.Duration
	shutdownTimeout time.Duration
	connContext     func(context.Context, net.Conn) context.Context
}

func defaultHTTPConfig() httpConfig {
	return httpConfig{
		shutdownTimeout: 5 * time.Second,
	}
}

func WithShutdownOnIdle(d time.Duration) HTTPOption {
	return func(c *httpConfig) {
		c.idleTimeout = d
	}
}

func WithShutdownTimeout(d time.Duration) HTTPOption {
	return func(c *httpConfig) {
		c.shutdownTimeout = d
	}
}

func WithConnContext(
	fn func(context.Context, net.Conn) context.Context,
) HTTPOption {
	return func(c *httpConfig) {
		c.connContext = fn
	}
}
