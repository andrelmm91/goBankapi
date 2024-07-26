package db

import (
	"context"
	"log"
	"os"
	"simplebank/util"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..") // app.env is on the root
	if err != nil {
		log.Fatal("cannot load configuration:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the DB:", err)
	}

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
