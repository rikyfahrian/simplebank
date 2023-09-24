package db

import (
	"database/sql"
	"log"
	"os"
	"techschool/util"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	config, err := util.LoadConfig("../../.env")
	if err != nil {
		panic(err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}

	testQueries = New(testDB)

	log.Println("connect to pg")

	os.Exit(m.Run())

}
