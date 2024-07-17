package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/util"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestVerifyEmailAPI(t *testing.T) {
	user, _ := randomUser(t)
	verifyEmail := randomVerifyEmail(user)

	type Query struct {
		emailId    int64
		secretCode string
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				emailId:    verifyEmail.ID,
				secretCode: verifyEmail.SecretCode,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.VerifyEmailTxParams{
					EmailId:    verifyEmail.ID,
					SecretCode: verifyEmail.SecretCode,
				}
				store.EXPECT().
					VerifyEmailTx(gomock.Any(), arg).
					Times(1).
					Return(db.VerifyEmailTxResult{User: user, VerifyEmail: verifyEmail}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				// requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "BadRequestFromQuery",
			query: Query{
				emailId:    verifyEmail.ID,
				secretCode: "",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					VerifyEmailTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			query: Query{
				emailId:    verifyEmail.ID,
				secretCode: verifyEmail.SecretCode,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.VerifyEmailTxParams{
					EmailId:    verifyEmail.ID,
					SecretCode: verifyEmail.SecretCode,
				}
				store.EXPECT().
					VerifyEmailTx(gomock.Any(), arg).
					Times(1).
					Return(db.VerifyEmailTxResult{}, fmt.Errorf("internal server error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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
			recorder := httptest.NewRecorder()

			url := "/verify_email"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("email_id", fmt.Sprintf("%d", tc.query.emailId))
			q.Add("secret_code", fmt.Sprintf("%s", tc.query.secretCode))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomVerifyEmail(user db.User) db.VerifyEmail {
	return db.VerifyEmail{
		ID:         util.RandomInt(1, 1000),
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
		IsUsed:     false,
		CreatedAt:  time.Now(),
		ExpiredAt:  time.Now().Add(10 * time.Second),
	}
}
