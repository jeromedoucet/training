package dao

import (
	"database/sql"
	"log"

	"github.com/jeromedoucet/training/configuration"
	_ "github.com/lib/pq"
)

// Conn handle db connection and DAO
// instanciation.
type Conn struct {
	db      *sql.DB
	UserDAO *UserDAO
}

// Open init connection to database and init all prepared statement.
// If something wrong happened, it panics.
func Open(conf *configuration.GlobalConf) (conn *Conn) {
	db, err := sql.Open("postgres", conf.DbStringConnection())
	if err != nil {
		log.Fatalf("An error is returned during db connection %s", err.Error())
	}

	res := &Conn{db: db}
	res.UserDAO = createUserDAO(res.db)

	return res
}
