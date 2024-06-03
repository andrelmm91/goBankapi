- create the migration files
  migrate create -ext sql -dir /migration -seq init_schema

- migrate db
  migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
