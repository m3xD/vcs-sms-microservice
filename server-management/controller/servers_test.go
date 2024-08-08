package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	mock_sqlc "server-management/db/mock"
	db "server-management/db/sqlc"
	"server-management/service"
	"server-management/token"
	"server-management/util"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
	"go.uber.org/mock/gomock"
)

func TestCreateServerAPI(t *testing.T) {
	server := randomServer()
	server.ID = 2
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
				"name":   server.Name,
				"status": server.Status,
				"ipv4":   server.Ipv4,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().CreateServer(gomock.Any(), gomock.Any()).Times(1).Return(server, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatch(t, recorder.Body, server)
			},
		},
		{
			nameTestCase: "bad request",
			body:         gin.H{},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().CreateServer(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			nameTestCase: "internal error",
			body: gin.H{
				"name":   server.Name,
				"status": server.Status,
				"ipv4":   server.Ipv4,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().CreateServer(gomock.Any(), gomock.Any()).Times(1).Return(db.Server{}, sql.ErrConnDone)
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

			url := "/create_server"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			require.NoError(t, err)

			tc.setupAuth(t, request, s.Token)
			s.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetServerAPI(t *testing.T) {
	server := make([]db.Server, 1)
	server[0] = randomServer()

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
				"page_id":   1,
				"page_size": 1,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetServer(gomock.Any(), gomock.Any()).Times(1).Return(server, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireArrayBodyMatch(t, recorder.Body, server)
			},
		},
		{
			nameTestCase: "bad request",
			body:         gin.H{},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetServer(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			nameTestCase: "internal error",
			body: gin.H{
				"page_id":   1,
				"page_size": 1,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetServer(gomock.Any(), gomock.Any()).Times(1).Return([]db.Server{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			nameTestCase: "not found",
			body: gin.H{
				"page_id":   1,
				"page_size": 1,
				"name":      "VCS",
				"isName":    true,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetServer(gomock.Any(), gomock.Any()).Times(1).Return([]db.Server{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			nameTestCase: "invalid id",
			body: gin.H{
				"page_id":   1,
				"page_size": 1,
				"id":        0,
				"isID":      true,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().UpdateServer(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			nameTestCase: "invalid ipv4",
			body: gin.H{
				"page_id":   1,
				"page_size": 1,
				"ipv4":      "",
				"isIpv4":    true,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().UpdateServer(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			nameTestCase: "invalid name",
			body: gin.H{
				"page_id":   1,
				"page_size": 1,
				"name":      "",
				"isName":    true,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().UpdateServer(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
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

			url := "/get_server"
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, s.Token)
			s.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateServerAPI(t *testing.T) {
	server := randomServer()
	server.ID = 1
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
				"id":      1,
				"name":    "VCS",
				"setName": true,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().UpdateServer(gomock.Any(), gomock.Any()).Times(1).Return(server, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatch(t, recorder.Body, server)
			},
		},
		{
			nameTestCase: "bad request",
			body:         gin.H{},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().UpdateServer(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			nameTestCase: "internal error",
			body: gin.H{
				"id":      1,
				"name":    "VCS",
				"setName": true,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().UpdateServer(gomock.Any(), gomock.Any()).Times(1).Return(db.Server{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			nameTestCase: "invalid name",
			body: gin.H{
				"id":      1,
				"name":    "",
				"setName": true,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().UpdateServer(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			nameTestCase: "invalid ipv4",
			body: gin.H{
				"id":      1,
				"ipv4":    "",
				"setIpv4": true,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().UpdateServer(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
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

			url := "/update_server"
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, s.Token)
			s.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteServerAPI(t *testing.T) {
	server := randomServer()
	server.ID = 1
	testCase := []struct {
		nameTestCase  string
		serverID      int
		buildStubs    func(store *mock_sqlc.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			nameTestCase: "OK",
			serverID:     int(server.ID),
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().DeleteServer(gomock.Any(), gomock.Eq(server.ID)).Times(1).Return(nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			nameTestCase: "bad request",
			serverID:     0,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().DeleteServer(gomock.Any(), gomock.Eq(server.ID)).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			nameTestCase: "internal error",
			serverID:     int(server.ID),
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().DeleteServer(gomock.Any(), gomock.Eq(server.ID)).Times(1).Return(sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
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

			url := fmt.Sprintf("/delete_server/%d", tc.serverID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, s.Token)
			s.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestImportExcelAPI(t *testing.T) {
	server := []db.Server{randomServer(), randomServer()}
	path, err := service.ExportServer(server)
	require.NoError(t, err)
	f, err := excelize.OpenFile(path)

	require.NoError(t, err)

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Errorf("%w", err)
		}
	}()

	tmpS := randomServer()

	f.SetCellValue("Result 1", "A4", 1)
	f.SetCellValue("Result 1", "B4", tmpS.Name)
	f.SetCellValue("Result 1", "C4", tmpS.Ipv4)
	f.SetCellValue("Result 1", "D4", tmpS.Status)
	f.SetCellValue("Result 1", "E4", time.Now())
	f.SetCellValue("Result 1", "F4", time.Now())

	f.Save()

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
				"path": path,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetAllServers(gomock.Any()).Times(1).Return(server, nil)
				store.EXPECT().CreateServer(gomock.Any(), gomock.Any()).Times(1).Return(tmpS, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			nameTestCase: "internal error",
			body: gin.H{
				path: "",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetAllServers(gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			nameTestCase: "error when create server",
			body: gin.H{
				"path": path,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetAllServers(gomock.Any()).Times(1).Return(server, nil)
				store.EXPECT().CreateServer(gomock.Any(), gomock.Any()).Times(1).Return(db.Server{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			nameTestCase: "error when get all server",
			body: gin.H{
				"path": path,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetAllServers(gomock.Any()).Times(1).Return([]db.Server{}, sql.ErrConnDone)
				store.EXPECT().CreateServer(gomock.Any(), gomock.Any()).Times(0)
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

			url := "/import_excel"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			require.NoError(t, err)

			tc.setupAuth(t, request, s.Token)
			s.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestExportExcelAPI(t *testing.T) {
	server := make([]db.Server, 1)
	server[0] = randomServer()

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
				"page_id":   1,
				"page_size": 1,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetServer(gomock.Any(), gomock.Any()).Times(1).Return(server, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			nameTestCase: "bad request",
			body:         gin.H{},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetServer(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			nameTestCase: "internal error",
			body: gin.H{
				"page_id":   1,
				"page_size": 1,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetServer(gomock.Any(), gomock.Any()).Times(1).Return([]db.Server{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			nameTestCase: "not found",
			body: gin.H{
				"page_id":   1,
				"page_size": 1,
				"name":      "VCS",
				"isName":    true,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetServer(gomock.Any(), gomock.Any()).Times(1).Return([]db.Server{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			nameTestCase: "invalid id",
			body: gin.H{
				"page_id":   1,
				"page_size": 1,
				"id":        0,
				"isID":      true,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().UpdateServer(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			nameTestCase: "invalid ipv4",
			body: gin.H{
				"page_id":   1,
				"page_size": 1,
				"ipv4":      "",
				"isIpv4":    true,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().UpdateServer(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			nameTestCase: "invalid name",
			body: gin.H{
				"page_id":   1,
				"page_size": 1,
				"name":      "",
				"isName":    true,
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().UpdateServer(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, request, tokenMaker, "bearer", username, scopes, role, 10*time.Minute)
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

			url := "/export_excel"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, s.Token)
			s.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomServer() db.Server {
	return db.Server{
		Name:   util.RandomServerName(),
		Status: int32(util.RandomStatus()),
		Ipv4:   util.RandomIP(),
	}
}

func requireBodyMatch(t *testing.T, body *bytes.Buffer, server db.Server) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var s db.Server
	err = json.Unmarshal(data, &s)
	require.NoError(t, err)
	require.Equal(t, server, s)
}

func requireArrayBodyMatch(t *testing.T, body *bytes.Buffer, server []db.Server) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var s []db.Server
	err = json.Unmarshal(data, &s)
	require.NoError(t, err)
	require.Equal(t, server, s)
}
