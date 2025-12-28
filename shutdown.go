package ctxsrv

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

func shutdownServer(ctx context.Context, srv *http.Server, logger *slog.Logger) {
	logger.Warn("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		if err := srv.Close(); err != nil {
			logger.Error("forced close failed", "error", err)
		}
	}
}
