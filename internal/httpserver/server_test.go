package httpserver_test

import (
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/gsh-assessment/internal/httpserver"
	"github.com/vcraescu/gsh-assessment/pkg/log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	t.Parallel()

	logger := log.NewLogger()
	srv := httpserver.New(logger)

	srv.Post("/orders", func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)

		return nil
	})
	srv.Get("/healthz", func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)

		return nil
	})

	addr, client, tearDown := setupTest(t, srv)

	t.Cleanup(tearDown)

	t.Run("method not allowed", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest(http.MethodGet, addr+"/orders", http.NoBody)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("success post", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest(http.MethodPost, addr+"/orders", http.NoBody)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("success post", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest(http.MethodGet, addr+"/healthz", http.NoBody)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func setupTest(t *testing.T, srv httpserver.Server) (address string, client *http.Client, tearDown func()) {
	t.Helper()

	testSrv := httptest.NewUnstartedServer(srv)
	testSrv.Start()

	return testSrv.URL, testSrv.Client(), func() {
		testSrv.Close()
	}
}
