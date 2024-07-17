- db migration
  github.com/golang-migrate/migrate/v4
  migrate create -ext sql -dir db/migration -seq add_sessions

- sqlc (manage sql crud)
  go get github.com/kyleconroy/sqlc/cmd/sqlc
  choco install kyleconroy/sqlc/sqlc
  sudo snap install sqlc

- codespace sudo (optional if using git codespace)
  sudo apt update
  sudo apt install snapd
  docker pull sqlc/sqlc
  docker run --rm -v $(pwd):/src -w /src sqlc/sqlc init
  docker run --rm -v $(pwd):/src -w /src sqlc/sqlc generate

- test sqlc DB
  go get github.com/lib/pq

- testing package GO
  go get -u github.com/stretchr/testify

- http framework GIN
  go get -u github.com/gin-gonic/gin

- go Viper for load config
  go get github.com/spf13/viper

- goMock for mocking DB testing
  go get github.com/golang/mock/mockgen@v1.6.0
  sudo apt install mockgen
  running mockgen in terminal:
  mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store

- Custom validator for go API
  go get github.com/go-playground/validator/v10

- Postgres driver and toolkit
  github.com/jackc/pgx/v5
  github.com/jackc/pgx/v5/pgconn

- generate ID to identify the token
  go get github.com/google/uuid

- JWT for Go token
  go get -u github.com/golang-jwt/jwt/v5

- PASETO for token mngt
  go get github.com/o1egl/paseto

- database migration
  update migration
  generate sqlc
  update mock

- message broker with ASYNQ and backed up with Redis
  go get -u github.com/hibiken/asynq

- logger JSON with high contrast
  go get -u github.com/rs/zerolog/log

- library for email (send, template)
  go get github.com/jordan-wright/email
