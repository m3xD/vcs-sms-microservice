package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	mock_sqlc "server-management/db/mock"
	"server-management/token"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateScope(t *testing.T) {
	testCase := []struct {
		nameTestCase  string
		body          gin.H
		buildStubs    func(store *mock_sqlc.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			nameTestCase: "OK",
			body: gin.H{
				"name": "api-test:scope_demo",
				"role": "role_test",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().
					CreateScope(gomock.Any(), gomock.Any()).
					Times(1)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.nameTestCase, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			store := mock_sqlc.NewMockStore(control)
			tc.buildStubs(store)

			s := NewTestServer(t, store)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/get_scope"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, s.Token)
			s.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		},
		)
	}
}
