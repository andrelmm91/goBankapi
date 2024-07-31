package api

import (
	"errors"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"

	"github.com/gin-gonic/gin"
)

// CreateAccountRequest represents the request body for creating a new account
type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"` // binding:currency is the custom Validator from validator.go
}

// createAccount creates a new account for the authenticated user
// @Summary Create a new account
// @Description Create a new account for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Param body body CreateAccountRequest true "Create account request body"
// @Success 200 {object} Account "Account created successfully"
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 500 {object} gin.H "Internal server error"
// @Security BearerAuth
// @Router /accounts [post]
func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest

	// validating the request from the body json.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get auth payload
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// getAccountRequest represents the URI parameters for getting an account
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getAccount retrieves account information by ID
// @Summary Get account by ID
// @Description Get details of an account by its ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} Account "Account details"
// @Failure 400 {object} gin.H "Invalid URI parameter"
// @Failure 404 {object} gin.H "Account not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Security BearerAuth
// @Router /accounts/{id} [get]
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	// validating the request from the URI.
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get auth payload
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// listAccountRequest represents the query parameters for listing accounts
type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5"`
}

// listAccounts lists all accounts for the authenticated user with pagination
// @Summary List accounts
// @Description List all accounts for the authenticated user with pagination
// @Tags accounts
// @Accept json
// @Produce json
// @Param page_id query int true "Page number"
// @Param page_size query int true "Page size"
// @Success 200 {array} Account "List of accounts"
// @Failure 400 {object} gin.H "Invalid query parameters"
// @Failure 500 {object} gin.H "Internal server error"
// @Security BearerAuth
// @Router /accounts [get]
func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get auth payload
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
