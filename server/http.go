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
	handler http.Handler,
	log *slog.Logger,
	ln net.Listener,
	idleTimeout time.Duration,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	srv := &http.Server{
		Handler:     handler,
		BaseContext: func(ln net.Listener) context.Context { return ctx },

		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      0, // MUST be 0 for streaming
		IdleTimeout:       0,
		ConnState:         idleTracker(ctx, cancel, idleTimeout),

		ErrorLog: slog.NewLogLogger(log.Handler(), slog.LevelError),
	}
	defer srv.Close()

	go func() {
		log.Info("HTTP server listening",
			slog.String("addr", ln.Addr().String()),
		)
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server failed",
				"error", err,
			)
			cancel()
		}
	}()

	<-ctx.Done()

	log.Warn("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Warn("server successfully shut down")

	return nil
}
