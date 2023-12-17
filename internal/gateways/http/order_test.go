package http_test

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/gsh-assessment/internal/adapters"
	"github.com/vcraescu/gsh-assessment/internal/domain"
	"github.com/vcraescu/gsh-assessment/internal/gateways/http"
	"github.com/vcraescu/gsh-assessment/pkg/echomw"
	"github.com/vcraescu/gsh-assessment/pkg/log"
	"io"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
)

type Context struct {
	echo.Context

	request *stdhttp.Request
}

func (c *Context) Request() *stdhttp.Request {
	return c.request
}

func TestNewCreateOrderHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		c echo.Context
	}

	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantBody       []byte
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "invalid request",
			args: args{
				c: &Context{
					request: newRequest(t, "test"),
				},
			},
			wantStatusCode: stdhttp.StatusBadRequest,
			wantBody:       marshalJSON(t, echomw.ErrorResponse{Error: domain.ErrInvalidArgument.Error()}),
		},
		{
			name: "empty request",
			args: args{
				c: &Context{
					request: newRequest(t, nil),
				},
			},
			wantStatusCode: stdhttp.StatusBadRequest,
			wantBody: marshalJSON(t, echomw.ErrorResponse{
				Error: "quantity must be greater than zero; got 0: invalid argument",
			}),
		},
		{
			name: "success",
			args: args{
				c: &Context{
					request: newRequest(t, http.CreateOrderRequest{Quantity: 251}),
				},
			},
			wantStatusCode: stdhttp.StatusOK,
			wantBody: marshalJSON(t, http.CreateOrderResponse{
				Data: domain.Order{
					Rows: []domain.OrderRow{
						{
							Quantity: 1,
							Pack:     500,
						},
					},
				},
			}),
		},
	}

	logger := log.NewNopLogger()

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, err := adapters.NewPackRepository()
			require.NoError(t, err)

			svc := domain.NewOrderService(repo)
			h := http.NewCreateOrderHandler(svc, logger)

			rec := httptest.NewRecorder()
			err = h(tt.args.c)

			if tt.wantErr != nil {
				tt.wantErr(t, err)

				return
			}

			require.NoError(t, err)

			got := rec.Result()

			require.Equal(t, tt.wantStatusCode, got.StatusCode)
			require.Equal(t, string(tt.wantBody), string(readBody(t, got)))
		})
	}
}

func newRequest(t *testing.T, body any) *stdhttp.Request {
	t.Helper()

	buf := &bytes.Buffer{}

	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		require.NoError(t, err)
	}

	return httptest.NewRequest(stdhttp.MethodGet, "http://example.com/test", buf)
}

func readBody(t *testing.T, resp *stdhttp.Response) []byte {
	t.Helper()

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return bytes.TrimSpace(b)
}

func marshalJSON(t *testing.T, v any) []byte {
	t.Helper()

	b, err := json.Marshal(v)
	require.NoError(t, err)

	return b
}
