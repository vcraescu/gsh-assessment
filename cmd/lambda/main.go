package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/labstack/echo/v4"
	"github.com/vcraescu/gsh-assessment/internal/adapters"
	"github.com/vcraescu/gsh-assessment/internal/domain"
	"github.com/vcraescu/gsh-assessment/internal/gateways/http"
	"github.com/vcraescu/gsh-assessment/pkg/echomw"
	"github.com/vcraescu/gsh-assessment/pkg/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"io"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
)

var (
	httpServer *httptest.Server
	tp         = trace.NewTracerProvider(trace.WithSampler(trace.AlwaysSample()))
	tracer     = tp.Tracer("lambda")
)

func main() {
	otel.SetTracerProvider(tp)

	logger := log.NewLogger()

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

	httpServer = httptest.NewServer(e)
	defer httpServer.Close()

	lambda.Start(handler)
}

func handler(ctx context.Context, in events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	ctx, span := tracer.Start(ctx, "handler")
	defer span.End()

	out := events.APIGatewayV2HTTPResponse{}

	req, err := apiGWRequestToHTTPRequest(ctx, in)
	if err != nil {
		return out, fmt.Errorf("apiGWRequestToHTTPRequest: %w", err)
	}
	defer req.Body.Close()

	resp, err := httpServer.Client().Do(req)
	if err != nil {
		return out, fmt.Errorf("do: %w", err)
	}

	return httpResponseToAPIGWResponse(resp)
}

func apiGWRequestToHTTPRequest(ctx context.Context, in events.APIGatewayV2HTTPRequest) (*stdhttp.Request, error) {
	rawURL := httpServer.URL + in.RawPath

	req, err := stdhttp.NewRequestWithContext(ctx, in.RequestContext.HTTP.Method, rawURL, strings.NewReader(in.Body))
	if err != nil {
		return nil, fmt.Errorf("newRequestWithContext: %w", err)
	}

	return req, nil
}

func httpResponseToAPIGWResponse(resp *stdhttp.Response) (events.APIGatewayV2HTTPResponse, error) {
	out := events.APIGatewayV2HTTPResponse{
		StatusCode:        resp.StatusCode,
		MultiValueHeaders: resp.Header,
		Headers:           make(map[string]string),
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, fmt.Errorf("readAll: %w", err)
	}

	for key, values := range resp.Header {
		out.Headers[key] = values[0]
	}

	out.Body = string(body)

	return out, nil
}
