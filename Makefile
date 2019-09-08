

db-up:
	sql-migrate up
	sql-migrate up --env=test

backend-test: db-up
	go test -p 1 ./...
