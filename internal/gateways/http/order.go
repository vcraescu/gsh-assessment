package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/vcraescu/gsh-assessment/internal/domain"
	"github.com/vcraescu/gsh-assessment/pkg/echomw"
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

func NewCreateOrderHandler(svc OrderService, logger log.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := CreateOrderRequest{}
		if err := c.Bind(&req); err != nil {
			return fmt.Errorf("bind: %w", err)
		}

		ctx := c.Request().Context()

		order, err := svc.Create(ctx, req.Quantity)
		if err != nil {
			logger.Error(ctx, "create order failed", slog.Any("payload", req), log.Error(err))

			return handleError(c, err)
		}

		return c.JSON(http.StatusOK, CreateOrderResponse{Data: order})
	}
}

func handleError(c echo.Context, err error) error {
	resp := echomw.ErrorResponse{
		Error: err.Error(),
	}

	code := http.StatusInternalServerError

	if errors.Is(err, domain.ErrInvalidArgument) {
		code = http.StatusBadRequest
	}

	return c.JSON(code, resp)
}
