package test

import (
	"database/sql"
	"log"

	"github.com/jeromedoucet/training/configuration"
)

func Conf() *configuration.GlobalConf {
	return &configuration.GlobalConf{DbName: "training_test", User: "postgres"}
}

// Clean empty all table in the training model
func CleanDB(c *configuration.GlobalConf) {
	db, err := sql.Open("postgres", c.DbStringConnection())
	if err != nil {
		log.Fatalf("An error is returned during db connection %s", err.Error())
	}

	_, err = db.Exec(`TRUNCATE "user"`)

	if err != nil {
		db.Close()
		log.Fatalf("An error is returned during db clean up %s", err.Error())
	}

	db.Close()

}
