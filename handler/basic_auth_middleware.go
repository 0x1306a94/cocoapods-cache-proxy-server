package handler

import (
	"cocoapods-cache-proxy-server/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BasicAuthMiddleware(authConfig *config.AuthorizationConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		if !authConfig.ValidationForBasicAuthorization(auth) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		ctx.Next()
	}
}
