package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nieyanzhai/simple-bank/util"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %s", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
