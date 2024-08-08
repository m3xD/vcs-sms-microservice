package middleware_test

import (
	"github.com/stretchr/testify/require"
	"testing"
	"vcssms/api"
	db "vcssms/db/sqlc"
)

func NewTestServer(t *testing.T, store db.Store) *api.Server {
	server, err := api.NewServer(store)
	require.NoError(t, err)

	return server
}
