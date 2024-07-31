package api

import (
	"errors"
	"fmt"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"

	"github.com/gin-gonic/gin"
)

// TransferRequest defines the request body for creating a transfer
// @Description Request body for initiating a transfer between accounts
// @Param from_account_id body int64 true "ID of the account to transfer from"
// @Param to_account_id body int64 true "ID of the account to transfer to"
// @Param amount body int64 true "Amount to transfer"
// @Param currency body string true "Currency of the accounts"
// @Accept json
// @Produce json
type TransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=1"`
	Currency      string `json:"currency" binding:"required,currency"` // binding:currency is the custom Validator from validator.go
}

// createTransfer handles the creation of a new transfer
// @Summary Create a Transfer
// @Description Initiate a transfer between two accounts. The request should include the account IDs, amount, and currency.
// @Tags transfers
// @Accept json
// @Produce json
// @Param request body TransferRequest true "Transfer Request"
// @Success 200 {object} db.TransferTxResult "Transfer successfully processed"
// @Failure 400 {object} gin.H "Bad Request - Invalid request data"
// @Failure 401 {object} gin.H "Unauthorized - User is not authorized for this transfer"
// @Failure 404 {object} gin.H "Not Found - Account not found"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /transfers [post]
func (server *Server) createTransfer(ctx *gin.Context) {
	var req TransferRequest

	// Validating the request from the body JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validating account and currency
	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	// Validating currency with the sender and receiver
	if !valid {
		return
	}

	// Get auth payload
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	// Validate RBAC
	err := RBAC(ctx, authPayload.Role, []string{util.DepositorRole})
	if err != nil {
		return
	}

	// Preparing to start the transaction
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// validAccount validates if an account exists and matches the provided currency
// @Description Check if the account exists and if its currency matches the provided currency
// @Param accountID query int64 true "Account ID to validate"
// @Param currency query string true "Currency to validate"
// @Success 200 {object} db.Account "Account successfully validated"
// @Failure 400 {object} gin.H "Bad Request - Currency mismatch"
// @Failure 404 {object} gin.H "Not Found - Account not found"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /accounts/validate [get]
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account (%d) currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
