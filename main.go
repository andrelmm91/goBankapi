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
	logz "github.com/rs/zerolog/log"
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

	// Connect to Redis and start task processor routine
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, store)

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

// REDIS queue
func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)

	logz.Info().Msg("start task processor")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("cannot start the task processor:", err)
	}
}