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