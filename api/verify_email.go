package api

import (
	"net/http"
	db "simplebank/db/sqlc"
	// db "simplebank/db/sqlc"
	// "simplebank/util"
	// "simplebank/worker"
	// "time"

	"github.com/gin-gonic/gin"
	// "github.com/hibiken/asynq"
)

type verifyEmailRequest struct {
	EmailId    int64  `form:"email_id" binding:"required"`
	SecretCode string `form:"secret_code" binding:"required"`
}

func (server *Server) verifyEmail(ctx *gin.Context) {
	var req verifyEmailRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// preparing to start the transaction
	arg := db.VerifyEmailTxParams{
		EmailId:    req.EmailId,
		SecretCode: req.SecretCode,
	}

	result, err := server.store.VerifyEmailTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}
