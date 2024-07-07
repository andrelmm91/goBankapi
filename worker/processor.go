package worker

import (
	"context"
	db "simplebank/db/sqlc"

	"github.com/hibiken/asynq"
)

type TaskProcess interface {
	start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error

}

type RedisTaskProcessor struct {
	server *asynq.Server
	store db.Store	
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcess {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{}, // empty means default configuration
	)

	return &RedisTaskProcessor{
		server: server,
		store: store,
	}
}

// similarly to http handler
func (processor *RedisTaskProcessor) start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}