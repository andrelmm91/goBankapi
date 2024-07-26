package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// store provides all functions to execute SQL db queries and transactions
type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

// NewStore create a new Store. connPool will allow many connection.
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		Queries: New(connPool),
		connPool: connPool,
	}
}


