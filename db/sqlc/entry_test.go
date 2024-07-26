package db

import (
	"context"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account1 Account) Entry {
	arg := CreateEntryParams{
		AccountID: account1.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testStore.CreateEntry(context.Background(), arg)

	// testing using package Testify
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.AccountID)
	require.NotZero(t, entry.Amount)

	return entry
}

// test CreateEntry
func TestCreateEntry(t *testing.T) {
	account1 := createRandomAccount(t)
	createRandomEntry(t, account1)
}

// test GetEntry
func TestGetEntry(t *testing.T) {
	account1 := createRandomAccount(t)

	entry1 := createRandomEntry(t, account1)
	entry2, err := testStore.GetEntry(context.Background(), entry1.ID)

	// testing using package Testify
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

// test ListEntries
func TestListEntries(t *testing.T) {
	account1 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, account1)
	}

	arg := ListEntriesParams{
		AccountID: account1.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testStore.ListEntries(context.Background(), arg)

	// testing using package Testify
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
