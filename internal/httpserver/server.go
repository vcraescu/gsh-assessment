package httpserver

import (
	"context"
	"github.com/vcraescu/gsh-assessment/pkg/log"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

type Server interface {
	http.Handler

	Post(pattern string, h HandlerFunc)
	Get(pattern string, h HandlerFunc)
}

type server struct {
	logger   log.Logger
	mux      *http.ServeMux
	handlers map[string]map[string]HandlerFunc
	once     sync.Once
}

func New(logger log.Logger) Server {
	return &server{
		logger:   logger,
		handlers: make(map[string]map[string]HandlerFunc),
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.setup(r.Context())

	s.mux.ServeHTTP(w, r)
}

func (s *server) Post(pattern string, h HandlerFunc) {
	s.registerHandlerFunc(http.MethodPost, pattern, h)
}

func (s *server) Get(pattern string, h HandlerFunc) {
	s.registerHandlerFunc(http.MethodGet, pattern, h)
}

func (s *server) setup(ctx context.Context) {
	s.once.Do(func() {
		s.mux = http.NewServeMux()

		for pattern, handlers := range s.handlers {
			s.mux.HandleFunc(pattern, s.newHandlerFunc(ctx, handlers))
		}
	})
}

func (s *server) registerHandlerFunc(method string, pattern string, h HandlerFunc) {
	if _, ok := s.handlers[pattern]; !ok {
		s.handlers[pattern] = make(map[string]HandlerFunc)
	}

	s.handlers[pattern][method] = h
}

func (s *server) newHandlerFunc(ctx context.Context, handlers map[string]HandlerFunc) http.HandlerFunc {
	return s.withLogger(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		h, ok := handlers[r.Method]
		if !ok {
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		s.handle(ctx, h)(w, r)
	})
}

func (s *server) handle(ctx context.Context, h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(ctx)

		if err := h(w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			s.logger.Error(ctx, "handler error", log.Error(err))

			return
		}
	}
}

func (s *server) withLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info(
			r.Context(), "REQUEST", slog.String("method", r.Method), slog.String("uri", r.RequestURI),
		)

		rec := httptest.NewRecorder()
		next(rec, r)

		for key, values := range rec.Header() {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(rec.Code)
		_, _ = io.Copy(w, rec.Body)

		s.logger.Info(r.Context(), "RESPONSE", slog.Int("code", rec.Code), slog.String("uri", r.RequestURI))
	}
}
