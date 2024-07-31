package api

import (
	"errors"
	"fmt"
	"net/http"
	"simplebank/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// authMiddleware is a middleware function that verifies the authorization token.
// It checks the "Authorization" header and validates the token.
// @Summary Authenticate API requests
// @Description Middleware to authenticate requests using Bearer tokens. Validates the token and sets the payload in the context.
// @Tags auth
// @Accept json
// @Produce json
// @Failure 401 {object} gin.H "Unauthorized"
// @Router / [get]  // This is a placeholder; actual routing does not apply to middleware.
// @Security BearerAuth
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			// abort the API call and return 401 to the user
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			// abort the API call and return 401 to the user
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			// abort the API call and return 401 to the user
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// payload is stored into a gin context with this specific key. To be used in the handlers.
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
