package httpx_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/gsh-assessment/internal/adapters"
	"github.com/vcraescu/gsh-assessment/internal/domain"
	"github.com/vcraescu/gsh-assessment/internal/gateways/httpx"
	"github.com/vcraescu/gsh-assessment/pkg/log"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewCreateOrderHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		req *http.Request
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
				req: newRequest(t, "test"),
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       marshalJSON(t, httpx.ErrorResponse{Error: domain.ErrInvalidArgument.Error()}),
		},
		{
			name: "empty request",
			args: args{
				req: newRequest(t, nil),
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody: marshalJSON(t, httpx.ErrorResponse{
				Error: "quantity must be greater than zero; got 0: invalid argument",
			}),
		},
		{
			name: "success",
			args: args{
				req: newRequest(t, httpx.CreateOrderRequest{Quantity: 251}),
			},
			wantStatusCode: http.StatusOK,
			wantBody: marshalJSON(t, httpx.CreateOrderResponse{
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
			h := httpx.NewCreateOrderHandler(svc, logger)

			rec := httptest.NewRecorder()
			err = h(rec, tt.args.req)
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

func newRequest(t *testing.T, body any) *http.Request {
	t.Helper()

	buf := &bytes.Buffer{}

	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		require.NoError(t, err)
	}

	return httptest.NewRequest(http.MethodGet, "http://example.com/test", buf)
}

func readBody(t *testing.T, resp *http.Response) []byte {
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
