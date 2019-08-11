package plan_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/jeromedoucet/training/configuration"
	"github.com/jeromedoucet/training/test"
)

var db *sql.DB
var conf *configuration.GlobalConf

func TestMain(m *testing.M) {
	var err error

	conf = test.Conf()
	conf.JwtSecret = "secret"
	db, err = sql.Open("postgres", conf.DbStringConnection())
	defer db.Close()
	if err != nil {
		log.Fatalf("An error is returned during db connection %s", err.Error())
	}

	os.Exit(m.Run())
}
