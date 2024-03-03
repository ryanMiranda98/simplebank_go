package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	util "github.com/ryanMiranda98/simplebank/util"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("Cannot read config:", err)
	}

	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	testQueries = New(testDb)
	os.Exit(m.Run())
}

func cleanUpDB() {
	testQueries.db.ExecContext(context.Background(), "DELETE FROM accounts;")
	testQueries.db.ExecContext(context.Background(), "DELETE FROM entries;")
	testQueries.db.ExecContext(context.Background(), "DELETE FROM transfers;")
	testQueries.db.ExecContext(context.Background(), "DELETE FROM users;")
}
