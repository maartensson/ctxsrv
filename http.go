package ctxsrv

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
)

func HTTP(
	ctx context.Context,
	port int,
	handler http.Handler,
	logger *slog.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	srv := &http.Server{
		Handler:     handler,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		ErrorLog:    slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	ln, err := listen(port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	go func() {
		logger.Info("HTTP server listening",
			slog.Int("port", port),
		)
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed",
				"error", err,
			)
			cancel()
		}
	}()

	<-ctx.Done()
	shutdownServer(context.Background(), srv, logger)
	return nil
}
