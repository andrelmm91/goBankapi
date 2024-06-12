package api

import (
	"net/http"
	"simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"require"`
	Currency string `json:"currency" binding:"require,oneof=USD EUR"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorRespose(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner: req.Owner,
		Currency: req.Currency,
		Balance: 0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorRespose(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}