package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RBAC(ctx *gin.Context, payloadRole string, accessibleRoles []string) error {

	if !hasPermission(payloadRole, accessibleRoles) {
		err := fmt.Errorf("permission denied")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return err
	}

	return nil
}

func hasPermission(userRole string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		if userRole == role {
			return true
		}
	}
	return false
}
