postgres:
	docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:15.1-alpine

createdb:
	docker exec -it postgres15 createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres15 dropdb --username=postgres simple_bank

migrateup:
	migrate -database "postgres://postgres:password@localhost:5433/simple_bank?sslmode=disable" -path db/migrations up

migratedown:
	migrate -database "postgres://postgres:password@localhost:5433/simple_bank?sslmode=disable" -path db/migrations down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test