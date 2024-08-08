package api

import (
	"fmt"
	"net/http"
	db "server-management/db/sqlc"
	"server-management/token"
	"server-management/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	username = "admin"
	role     = "admin"
	scopes   = []string{
		"api-server:read",
		"api-server:write",
		"api-user:read",
		"api-user:write",
		"api-report:read",
	}
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		SecretKey: "123",
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func AddAuthorization(t *testing.T, request *http.Request, tokenMarker token.Maker, authorizationType string, username string, scopes []string, role string, duration time.Duration) {
	token, err := tokenMarker.CreateToken(username, role, scopes, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	authorHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set("authorization", authorHeader)
}
