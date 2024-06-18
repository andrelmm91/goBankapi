- db migration
  github.com/golang-migrate/migrate/v4

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