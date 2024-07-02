package api

import (
	"bytes"
	"time"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "simplebank/db/mock"
	"simplebank/token"
	"simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestToken(t *testing.T) {
	user, _ := randomUser(t)
	_, mockRefreshPayload := randomRefreshToken(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"refresh_token": user.Username,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(mockRefreshPayload.ID)).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		// {
		// 	name: "UserNotFound",
		// 	body: gin.H{
		// 		"username": user.Username,
		// 		"password": password,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			GetUser(gomock.Any(), gomock.Any()).
		// 			Times(1).
		// 			Return(db.User{}, db.ErrRecordNotFound)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusNotFound, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "IncorrectPassword",
		// 	body: gin.H{
		// 		"username": user.Username,
		// 		"password": "incorrect",
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			GetUser(gomock.Any(), gomock.Eq(user.Username)).
		// 			Times(1).
		// 			Return(user, nil)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusUnauthorized, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "InternalError",
		// 	body: gin.H{
		// 		"username": user.Username,
		// 		"password": password,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			GetUser(gomock.Any(), gomock.Any()).
		// 			Times(1).
		// 			Return(db.User{}, sql.ErrConnDone)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "InvalidUsername",
		// 	body: gin.H{
		// 		"username": "invalid-user#1",
		// 		"password": password,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			GetUser(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomRefreshToken(t *testing.T) (mockRefreshToken string, payload *token.Payload) {
	maker, err := token.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	mockRefreshToken, payload, err = maker.CreateToken(util.RandomOwner(), time.Second)
	require.NotEmpty(t, mockRefreshToken)
	require.NotEmpty(t, payload)
	require.NoError(t, err)

	return
}
