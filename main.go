package main

import (
	"context"
	"log"

	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/mail"
	"simplebank/util"
	"simplebank/worker"

	"github.com/hibiken/asynq"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	logz "github.com/rs/zerolog/log"
)

func main() {
	config, err := util.LoadConfig(".") // . means that the path is the same as main.go
	if err != nil {
		log.Fatal("cannot load configurations:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the DB:", err)
	}

	store := db.NewStore(connPool)

	// Connect to Redis and start task processor routine
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(config, redisOpt, store)

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
func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)

	logz.Info().Msg("start task processor")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("cannot start the task processor:", err)
	}
}
