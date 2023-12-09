package httpserver

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type tracedServer struct {
	server Server
	tracer trace.Tracer
}

func NewTraced(server Server, tracer trace.Tracer) Server {
	return &tracedServer{
		server: server,
		tracer: tracer,
	}
}

func (t *tracedServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	otelhttp.NewHandler(t.server, "serverHTTP").ServeHTTP(w, r)
}

func (t *tracedServer) Post(pattern string, h HandlerFunc) {
	t.server.Post(pattern, h)
}

func (t *tracedServer) Get(pattern string, h HandlerFunc) {
	t.server.Get(pattern, h)
}
