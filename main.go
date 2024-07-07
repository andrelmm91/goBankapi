package main

import (
	"database/sql"
	"log"

	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"simplebank/worker"

	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".") // . means that the path is the same as main.go
	if err != nil {
		log.Fatal("cannot load configurations:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the DB:", err)
	}

	store := db.NewStore(conn)

	// Connect to Redis
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	// start http Server
	server, err := api.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal("cannot create the server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start the server:", err)
	}

}