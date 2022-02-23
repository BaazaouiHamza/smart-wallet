package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"git.digitus.me/pfe/smart-wallet/util"
	_ "github.com/lib/pq"
)

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://postgres:07719811gg@localhost:5432/pfe?sslmode=disable"
// )

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDb)

	os.Exit(m.Run())
}
