package middleware

import (
	"net/http"
	"server-management/token"

	"github.com/gin-gonic/gin"
)

func CheckScope(currentScope string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		scopes := ctx.MustGet("authorization_payload").(*token.Payload)
		for _, data := range scopes.Scope {
			if string(data) == currentScope {
				ctx.Next()
				return
			}
		}
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "not sufficient scope",
		})
		ctx.Abort()
		return
	}
}
