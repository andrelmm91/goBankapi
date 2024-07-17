createdb:
	docker exec -it db-postgres-1 createdb --username=root --owner=root simple_bank

dropdb: 
	docker exec -it db-postgres-1 dropdb simple_bank

migrateup:
	cd ./db/ && \
	migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	cd ./db/ && \
	migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	cd ./db/ && \
	sqlc generate

test:
	sudo go test -v -cover -short ./...

dcUpdate:
	docker-compose down
	docker image rm gobankapi-bank-backend
	sudo rm -rf db/postgres-data
	docker compose up

server: 
	go run main.go

mock:
	mockgen -package mockdb -destination ./db/mock/store.go simplebank/db/sqlc Store

.PHONY: createdb dropdb migratedown migrateup migratedown1 migrateup1 sqlc dcUpdate server mock new_migration