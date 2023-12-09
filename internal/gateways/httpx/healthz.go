package httpx

import (
	"github.com/vcraescu/gsh-assessment/internal/httpserver"
	"net/http"
)

func NewHealthzCheckHandler() httpserver.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)

		return nil
	}
}
