package api

import (
	"net/http"
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

// verifyEmailRequest defines the request parameters for verifying an email
// @Description Request parameters for email verification
// @Param email_id query int64 true "ID of the email to verify" example(12345)
// @Param secret_code query string true "Secret code for email verification" example("abcd1234")
// @Accept json
// @Produce json
type verifyEmailRequest struct {
	EmailId    int64  `form:"email_id" binding:"required"`
	SecretCode string `form:"secret_code" binding:"required"`
}

// verifyEmailResponse defines the response body for email verification
// @Description Response body indicating whether the email has been verified
// @Property is_verified boolean "Indicates if the email is verified" example(true)
// @Accept json
// @Produce json
type verifyEmailResponse struct {
	IsVerified bool `json:"is_verified"`
}

// verifyEmail handles the email verification process
// @Summary Verify Email
// @Description Verify an email address using the provided email ID and secret code. Returns whether the email was successfully verified.
// @Tags email
// @Accept json
// @Produce json
// @Param request query verifyEmailRequest true "Verify Email Request"
// @Success 200 {object} verifyEmailResponse "Email verified successfully"
// @Failure 400 {object} gin.H "Bad Request - Invalid input data"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /verify-email [get]
func (server *Server) verifyEmail(ctx *gin.Context) {
	var req verifyEmailRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Preparing to start the transaction
	arg := db.VerifyEmailTxParams{
		EmailId:    req.EmailId,
		SecretCode: req.SecretCode,
	}

	result, err := server.store.VerifyEmailTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := verifyEmailResponse{
		IsVerified: result.User.IsEmailVerified,
	}

	ctx.JSON(http.StatusOK, rsp)
}
