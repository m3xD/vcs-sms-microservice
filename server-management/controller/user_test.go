package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	mock_sqlc "server-management/db/mock"
	db "server-management/db/sqlc"
	"server-management/token"
	"server-management/util"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAPICreateUse(t *testing.T) {

	user := db.User{
		Username: "user_test",
		Password: "password",
		Email:    "example@gmail.com",
		Role:     "user",
	}

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
				"username": user.Username,
				"password": user.Password,
				"email":    user.Email,
				"role":     user.Role,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(user, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			nameTestCase: "error when hashing password",
			body: gin.H{
				"username": user.Username,
				"password": "321321d0psjklsadjlsakjdsakd908e12hjdksaljhdsakldjddsadsadsaddaasddshsadksadlkas",
				"email":    user.Email,
				"role":     user.Role,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			nameTestCase: "error when binding json",
			body:         gin.H{},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			nameTestCase: "error when create user",
			body: gin.H{
				"username": user.Username,
				"password": user.Password,
				"email":    user.Email,
				"role":     user.Role,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			require.NoError(t, err)

			tc.setupAuth(t, request, s.Token)
			s.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestAPILoginUse(t *testing.T) {
	pass, _ := util.HashPassword("user")
	user := db.User{
		Username: "user",
		Password: pass,
		Email:    "example@email.com",
		Role:     "user",
	}

	testCase := []struct {
		nameTestCase  string
		body          gin.H
		buildStubs    func(store *mock_sqlc.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			nameTestCase: "OK",
			body: gin.H{
				"username": "user",
				"password": "user",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().GetScope(gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			nameTestCase: "error when binding json",
			body:         gin.H{},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetScope(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			nameTestCase: "user not found",
			body: gin.H{
				"username": "userx",
				"password": "user",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrNoRows)
				store.EXPECT().GetScope(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			nameTestCase: "error when get user",
			body: gin.H{
				"username": "user",
				"password": "user",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)
				store.EXPECT().GetScope(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			nameTestCase: "error when compare password",
			body: gin.H{
				"username": "user",
				"password": "x",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().GetScope(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			nameTestCase: "error when get scope",
			body: gin.H{
				"username": "user",
				"password": "user",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().GetScope(gomock.Any(), gomock.Any()).Times(1).Return([]string{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := "/users/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			require.NoError(t, err)

			s.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
func TestAPIUpdateRole(t *testing.T) {
	pass, _ := util.HashPassword("user")
	user := db.User{
		ID:       1,
		Username: "user",
		Password: pass,
		Email:    "example@email.com",
		Role:     "user",
	}

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
				"id":   user.ID,
				"role": "admin",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				user.Role = "admin"
				store.EXPECT().UpdateRole(gomock.Any(), gomock.Any()).Times(1).Return(user, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			nameTestCase: "error when binding json",
			body:         gin.H{},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().UpdateRole(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			nameTestCase: "error when update user",
			body: gin.H{
				"id":   user.ID,
				"role": "admin",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().UpdateRole(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := "/update_role"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			require.NoError(t, err)

			tc.setupAuth(t, request, s.Token)
			s.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
