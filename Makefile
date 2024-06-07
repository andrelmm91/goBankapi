createdb:
	docker exec -it db-postgres-1 createdb --username=root --owner=root simple_bank

dropdb: 
	docker exec -it db-postgres-1 dropdb simple_bank

migrateup:
	migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	sudo go test -v -cover ./...

.PHONY: createdb dropdb migratedown migrateup sqlc