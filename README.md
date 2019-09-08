# Training application

Training is a tool that is designed to help building sport training plan.

# Backend

## setup

Training backend require [Golang](https://golang.org/) 11+.

Training use [postgres](https://www.postgresql.org/) as database. Two differents database are used `training` and `training_test` (used by integration tests). It must be created before going further.

Example for docker : 

```
$ docker run --name some-postgres -p 5432:5432 -d postgres
$ docker exec -ti --user postgres some-postgres bash -c "psql -c 'CREATE DATABASE training;'"
$ docker exec -ti --user postgres some-postgres bash -c "psql -c 'CREATE DATABASE training_test;'"
```

[sql-migrate](https://github.com/rubenv/sql-migrate) handle the migrations. To run migration, exec the following commands: 

```
sql-migrate up
sql-migrate up --env=test
```

or `make db-up`

The dependencies are managed by go module. You can load all of through a simple `go build`.

## tests

You can simply use `make backend-test` shortcut. If you want to run yourself the go command, be aware that many test are integration tests. This means that no parallel run is possible (a test specific dataset is loaded into the database). You may execute the test with `-p 1` option, to avoid package concurrent test execution.


# Frontend

The front root directory is ./front.

[yarn](https://yarnpkg.com/lang/en/) is the first thing to install. After that, a simple `yarn install` should install all dependencies.
