package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// VerifyEmailTxParams contains the input parameters of the VerifyEmail transaction
type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

// VerifyEmailTxResult is the result of the VerifyEmail transaction
type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

// VerifyEmailTx performs a monez VerifyEmail from one account to another
// It creates a VerifyEmail record, add account entries, update accounts balance within a single datybase transaction
func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})

		return err
	})

	return result, err
}
