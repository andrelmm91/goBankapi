package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRenewAccessToken(t *testing.T) {
	role := util.DepositorRole
	user, _ := randomUser(t, role)

	// Create a token maker instance with a fixed secret
	tokenMaker, err := token.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	// tokenMaker2, err := token.NewPasetoMaker(util.RandomString(32))

	mockRefreshToken, mockRefreshPayload := randomRefreshToken(t, tokenMaker, role, user.Username)
	// _, mockRefreshPayload2 := randomRefreshToken(t, tokenMaker2, util.RandomOwner())

	mockSession := randomSession(mockRefreshToken, mockRefreshPayload.Username, mockRefreshPayload, false, 1)
	mockSessionBlocked := randomSession(mockRefreshToken, mockRefreshPayload.Username, mockRefreshPayload, true, 1)
	mockSessionWrongUser := randomSession(mockRefreshToken, "anyUser", mockRefreshPayload, false, 1)
	mockSessionWrongRefreshToken := randomSession(util.RandomString(6), mockRefreshPayload.Username, mockRefreshPayload, false, 1)
	mockSessionExpired := randomSession(mockRefreshToken, mockRefreshPayload.Username, mockRefreshPayload, false, -1)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"refresh_token": mockRefreshToken,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(mockRefreshPayload.ID)).
					Times(1).
					Return(mockSession, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var rsp RenewAccessTokenResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &rsp)
				require.NoError(t, err)
				require.NotEmpty(t, rsp.AccessToken)
			},
		},
		{
			name: "BadRequestRequest",
			body: gin.H{
				"refresh_token": "",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UnauthorizedVerifiedToken",
			body: gin.H{
				"refresh_token": util.RandomString(6),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "SessionNotFound",
			body: gin.H{
				"refresh_token": mockRefreshToken,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{}, db.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "sessionBlocked",
			body: gin.H{
				"refresh_token": mockRefreshToken,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(mockRefreshPayload.ID)).
					Times(1).
					Return(mockSessionBlocked, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				// Parse the response body
				var responseBody map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				// Check the error message
				require.Equal(t, "blocked session", responseBody["error"])
			},
		},
		{
			name: "sessionDifferentUserName",
			body: gin.H{
				"refresh_token": mockRefreshToken,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(mockRefreshPayload.ID)).
					Times(1).
					Return(mockSessionWrongUser, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				// Parse the response body
				var responseBody map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				// Check the error message
				require.Equal(t, "incorrect session user", responseBody["error"])
			},
		},
		{
			name: "sessionDifferentRefreshToken",
			body: gin.H{
				"refresh_token": mockRefreshToken,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(mockRefreshPayload.ID)).
					Times(1).
					Return(mockSessionWrongRefreshToken, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				// Parse the response body
				var responseBody map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				// Check the error message
				require.Equal(t, "mismatched session token", responseBody["error"])
			},
		},
		{
			name: "sessionExpired",
			body: gin.H{
				"refresh_token": mockRefreshToken,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(mockRefreshPayload.ID)).
					Times(1).
					Return(mockSessionExpired, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				// Parse the response body
				var responseBody map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				// Check the error message
				require.Equal(t, "expired session", responseBody["error"])
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			server.tokenMaker = tokenMaker // Ensure the same token maker is used
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/tokens/renew_access"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomRefreshToken(t *testing.T, maker token.Maker, role string, username string) (string, *token.Payload) {
	mockRefreshToken, mockRefreshPayload, err := maker.CreateToken(username, role, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, mockRefreshToken)
	require.NotEmpty(t, mockRefreshPayload)

	return mockRefreshToken, mockRefreshPayload
}

func randomSession(mockRefreshToken string, username string, mockRefreshPayload *token.Payload, isBlocked bool, expiration time.Duration) db.Session {
	return db.Session{
		ID:           mockRefreshPayload.ID,
		Username:     username,
		RefreshToken: mockRefreshToken,
		UserAgent:    util.RandomOwner(),
		ClientIp:     util.RandomString(8),
		IsBlocked:    isBlocked,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(expiration * time.Minute),
	}
}
