package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

func HTTPS(
	ctx context.Context,
	port int,
	tlsConfig *tls.Config,
	router http.Handler,
	logger *slog.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	srv := &http.Server{
		Handler:     router,
		BaseContext: func(net.Listener) context.Context { return ctx },

		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      0, // MUST be 0 for streaming
		IdleTimeout:       0,

		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	ln, err := listenTLS(port, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	go func() {
		defer cancel()

		logger.Info("HTTPS server listening",
			"port", port,
		)

		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTPS server failed",
				"error",
				err,
			)
		}
	}()

	<-ctx.Done()

	shutdownServer(context.Background(), srv, logger)

	return nil
}
