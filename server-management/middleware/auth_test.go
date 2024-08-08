package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"vcssms/middleware"
	"vcssms/token"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func AddAuthorization(t *testing.T, request *http.Request, tokenMarker token.Maker, authorizationType string, username string, scopes []string, role string, duration time.Duration) {
	token, err := tokenMarker.CreateToken(username, role, scopes, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	authorHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set("authorization", authorHeader)
}

func TestAuthMiddleware(t *testing.T) {
	username := "user"
	role := "user"
	scopes := []string{
		"api-server:read",
		"api-user:read",
		"api-report:read",
	}

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "unsupported", username, scopes, role, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "", username, scopes, role, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "invalid username",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", "invalid username", scopes, role, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := NewTestServer(t, nil)
			authPath := "/auth"
			server.Router.GET(
				authPath,
				middleware.AuthMiddleWare(server.Token),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.Token)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
