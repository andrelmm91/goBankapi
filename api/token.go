package api

import (
	"errors"
	"fmt"
	"net/http"
	db "simplebank/db/sqlc"
	"time"

	"github.com/gin-gonic/gin"
)

// RenewAccessTokenRequest defines the request body for renewing access tokens
// @Description Request body for renewing access tokens
// @Param refresh_token body string true "Refresh Token"
// @Accept json
// @Produce json
type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RenewAccessTokenResponse defines the response body for renewing access tokens
// @Description Response body for renewing access tokens
// @Param access_token query string true "Access Token"
// @Param access_token_expires_at query string true "Access Token Expiration Time"
// @Accept json
// @Produce json
type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// RenewAccessToken handles the renewal of access tokens
// @Summary Renew Access Token
// @Description Renew an access token using a valid refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RenewAccessTokenRequest true "Renew Access Token Request"
// @Success 200 {object} RenewAccessTokenResponse "Access token successfully renewed"
// @Failure 400 {object} gin.H "Bad Request - Invalid refresh token"
// @Failure 401 {object} gin.H "Unauthorized - Invalid refresh token or session"
// @Failure 404 {object} gin.H "Not Found - Session not found"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /tokens/renew_access [post]
func (server *Server) RenewAccessToken(ctx *gin.Context) {
	var req RenewAccessTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verify the refresh token
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Find the session in DB
	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Check if the session is blocked
	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Check if session username matches the refresh token username
	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Check if session refresh token matches the request refresh token
	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Check expiration time
	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, refreshPayload.Role, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
