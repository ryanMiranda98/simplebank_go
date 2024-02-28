postgres:
	docker run --name pg-db -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres

createdb:
	docker exec -it pg-db createdb --username=postgres --owner=postgres simplebank

dropdb:
	docker exec -it pg-db dropdb --username=postgres simplebank

# create_migration_schema:
# 	migrate create -ext sql -dir db/migration -seq <schema_name>

migrate_up:
	migrate -database "postgres://postgres:mysecretpassword@localhost:5432/simplebank?sslmode=disable" -path db/migration up

migrate_down:
	migrate -database "postgres://postgres:mysecretpassword@localhost:5432/simplebank?sslmode=disable" -path db/migration down

sqlc_init:
	.\sqlc init

sqlc_generate:
	.\sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/ryanMiranda98/simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrate_up migrate_down sqlc_init sqlc_generate test server mock