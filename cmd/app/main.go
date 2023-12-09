package main

import (
	"context"
	_ "embed"
	"github.com/vcraescu/gsh-assessment/internal/adapters"
	"github.com/vcraescu/gsh-assessment/internal/domain"
	"github.com/vcraescu/gsh-assessment/internal/gateways/httpx"
	"github.com/vcraescu/gsh-assessment/internal/httpserver"
	"github.com/vcraescu/gsh-assessment/pkg/log"
	"github.com/vcraescu/gsh-assessment/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
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
		srv    = httpserver.NewTraced(httpserver.New(logger), tracer)
	)

	defer tp.Shutdown(ctx)

	repository, err := adapters.NewPackRepository()
	if err != nil {
		panic(err)
	}

	svc := domain.NewOrderService(repository)
	createOrderHandler := httpx.NewCreateOrderHandler(svc, logger)
	healthzCheckHandler := httpx.NewHealthzCheckHandler()

	srv.Get("/", web.StaticHandler)
	srv.Post("/orders", createOrderHandler)
	srv.Get("/healthz", healthzCheckHandler)

	if err := httpserver.Start(ctx, logger, srv, serverAddress); err != nil {
		panic(err)
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
