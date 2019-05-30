package configuration_test

import (
	"testing"

	"github.com/jeromedoucet/training/configuration"
)

func TestDbStringConnection(t *testing.T) {
	conf := configuration.GlobalConf{}

	if conf.DbStringConnection() != "" {
		t.Fatalf("expect %s, got %s", "", conf.DbStringConnection())
	}

	conf.DbName = "some_db_name"

	if conf.DbStringConnection() != "dbname=some_db_name" {
		t.Fatalf("expect %s, got %s", "dbname=some_db_name", conf.DbStringConnection())
	}

	conf.User = "some_user"

	if conf.DbStringConnection() != "dbname=some_db_name user=some_user" {
		t.Fatalf("expect %s, got %s", "dbname=some_db_name user=some_user", conf.DbStringConnection())
	}

	conf.Password = "some_password"

	if conf.DbStringConnection() != "dbname=some_db_name user=some_user password=some_password" {
		t.Fatalf("expect %s, got %s", "dbname=some_db_name user=some_user password=some_password", conf.DbStringConnection())
	}

	conf.Host = "localhost"

	if conf.DbStringConnection() != "dbname=some_db_name user=some_user password=some_password host=localhost" {
		t.Fatalf("expect %s, got %s", "dbname=some_db_name user=some_user password=some_password host=localhost", conf.DbStringConnection())
	}

	conf.Port = 5432

	if conf.DbStringConnection() != "dbname=some_db_name user=some_user password=some_password host=localhost port=5432" {
		t.Fatalf("expect %s, got %s", "dbname=some_db_name user=some_user password=some_password host=localhost port=5432", conf.DbStringConnection())
	}

}
