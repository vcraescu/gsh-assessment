package httpx

import (
	"context"
	"errors"
	"fmt"
	"github.com/vcraescu/gsh-assessment/internal/domain"
	"github.com/vcraescu/gsh-assessment/internal/httpserver"
	"github.com/vcraescu/gsh-assessment/pkg/log"
	"log/slog"
	"net/http"
)

type OrderService interface {
	Create(ctx context.Context, quantity int) (domain.Order, error)
}

type CreateOrderRequest struct {
	Quantity int `json:"quantity"`
}

type CreateOrderResponse struct {
	Data domain.Order `json:"data"`
}

func NewCreateOrderHandler(svc OrderService, logger log.Logger) httpserver.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		req := &CreateOrderRequest{}

		if err := decodeRequest(r, req); err != nil {
			return handleError(err, w)
		}

		order, err := svc.Create(r.Context(), req.Quantity)
		if err != nil {
			logger.Error(r.Context(), "create order failed", slog.Any("payload", req), log.Error(err))

			return handleError(err, w)
		}

		if err := encodeResponse(w, http.StatusOK, CreateOrderResponse{Data: order}); err != nil {
			return fmt.Errorf("encodeResponse: %w", err)
		}

		return nil
	}
}

func handleError(err error, w http.ResponseWriter) error {
	resp := ErrorResponse{
		Error: err.Error(),
	}

	code := http.StatusInternalServerError

	if errors.Is(err, domain.ErrInvalidArgument) {
		code = http.StatusBadRequest
	}

	if err := encodeResponse(w, code, resp); err != nil {
		return fmt.Errorf("encode error response: %w", err)
	}

	return nil
}
