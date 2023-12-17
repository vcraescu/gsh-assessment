package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/vcraescu/gsh-assessment/internal/adapters"
	"github.com/vcraescu/gsh-assessment/internal/domain"
	"github.com/vcraescu/gsh-assessment/internal/gateways/http"
	"github.com/vcraescu/gsh-assessment/pkg/echomw"
	"github.com/vcraescu/gsh-assessment/pkg/log"
	"github.com/vcraescu/gsh-assessment/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"log/slog"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	gracefulShutdownTimeout = time.Second
	serverAddress           = ":3000"
)

func main() {
	tp := trace.NewTracerProvider(trace.WithSampler(trace.AlwaysSample()))
	otel.SetTracerProvider(tp)

	var (
		tracer = otel.Tracer("app")
		logger = log.NewLogger()
		ctx    = gracefulShutdown(context.Background(), logger)
	)

	defer tp.Shutdown(ctx)

	repository, err := adapters.NewPackRepository()
	if err != nil {
		panic(err)
	}

	svc := domain.NewOrderService(repository)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(echomw.WithTracer(tracer), echomw.WithLogger(logger))

	e.GET("/healthz", http.NewHealthzHandler())
	e.POST("/orders", http.NewCreateOrderHandler(svc, logger))
	e.StaticFS("/", web.FS)

	go func() {
		<-ctx.Done()

		if err := e.Shutdown(ctx); err != nil {
			logger.Error(ctx, "shutdown failed", log.Error(err))
		}
	}()

	logger.Info(ctx, "server started", slog.String("address", serverAddress))

	if err := e.Start(serverAddress); err != nil {
		if !errors.Is(err, stdhttp.ErrServerClosed) {
			panic(fmt.Errorf("start: %w", err))
		}
	}
}

func gracefulShutdown(ctx context.Context, logger log.Logger) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer cancel()
		<-sigCh

		logger.Info(ctx, "shutting down server...")

		time.Sleep(gracefulShutdownTimeout)
	}()

	return ctx
}
