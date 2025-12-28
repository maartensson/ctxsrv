package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

func HTTP(
	ctx context.Context,
	cancel context.CancelFunc,
	handler http.Handler,
	log *slog.Logger,
	ln net.Listener,
	idleTimeout time.Duration,
	connContext func(ctx context.Context, c net.Conn) context.Context,
) error {

	srv := &http.Server{
		Handler:           handler,
		BaseContext:       func(ln net.Listener) context.Context { return ctx },
		ConnContext:       connContext,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      0, // MUST be 0 for streaming
		IdleTimeout:       0,
		ConnState:         idleTracker(ctx, cancel, idleTimeout),
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

	shutdownCtx, shutDownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutDownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Warn("server successfully shut down")

	return nil
}
