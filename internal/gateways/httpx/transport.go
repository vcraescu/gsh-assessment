package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vcraescu/gsh-assessment/internal/domain"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

func decodeRequest(r *http.Request, request any) error {
	if r.Body == nil {
		return nil
	}

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		if !errors.Is(err, io.EOF) {
			return domain.ErrInvalidArgument
		}
	}

	return nil
}

func encodeResponse(w http.ResponseWriter, code int, resp any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	return nil
}
