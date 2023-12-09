package adapters_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/gsh-assessment/internal/adapters"
	"testing"
)

func TestPackRepository_FindAll(t *testing.T) {
	t.Parallel()

	r, err := adapters.NewPackRepository()
	require.NoError(t, err)

	got, err := r.FindAll(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, got)
}
