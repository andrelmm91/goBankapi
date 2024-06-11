createdb:
	docker exec -it db-postgres-1 createdb --username=root --owner=root simple_bank

dropdb: 
	docker exec -it db-postgres-1 dropdb simple_bank

migrateup:
	cd ./db/ && \
	migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	cd ./db/ && \
	sqlc generate

test:
	sudo go test -v -cover ./...

setupDBviaDockerCompose:
	cd ./db/ && \
	docker-compose up -d

.PHONY: createdb dropdb migratedown migrateup sqlc setupDBviaDockerCompose