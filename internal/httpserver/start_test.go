package httpserver_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/gsh-assessment/internal/httpserver"
	"github.com/vcraescu/gsh-assessment/pkg/log"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	t.Parallel()

	var (
		srv         = httpserver.New(log.NewLogger())
		logger      = log.NewNopLogger()
		ctx, cancel = context.WithCancel(context.Background())
		done        = make(chan struct{})
	)

	go func() {
		defer close(done)

		err := httpserver.Start(ctx, logger, srv, ":54666")
		require.NoError(t, err)
	}()

	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		require.Fail(t, "server didn't shutdown")
	}
}
