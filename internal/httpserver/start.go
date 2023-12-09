package httpserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/vcraescu/gsh-assessment/pkg/log"
	"log/slog"
	"net/http"
)

func Start(ctx context.Context, logger log.Logger, srv Server, address string) error {
	httpSrv := http.Server{
		Addr: address,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(ctx)

			srv.ServeHTTP(w, r)
		}),
	}

	go func() {
		<-ctx.Done()
		if err := httpSrv.Shutdown(ctx); err != nil {
			logger.Error(ctx, "shutdown failed", log.Error(err))
		}
	}()

	logger.Info(ctx, "server started", slog.String("address", address))

	if err := httpSrv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listenAndServe: %w", err)
		}
	}

	return nil
}
