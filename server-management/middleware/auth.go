package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	db "server-management/db/sqlc"
	"server-management/token"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func AuthMiddleWare(tokenMarker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("authorization")
		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is not provided"})
			return
		}

		field := strings.Fields(authorizationHeader)
		if len(field) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}
		authorizationType := strings.ToLower(field[0])
		if authorizationType != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("unsupported authorization type %s", authorizationType)})
		}

		accessToken := field[1]

		payload, err := tokenMarker.VerifyToken(accessToken)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c := context.Background()
		conn, err := pgx.Connect(c, os.Getenv("DB_SOURCE"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer conn.Close(c)

		store := db.NewStore(conn)

		u, err := store.GetUser(ctx, payload.Username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		scopes, err := store.GetScope(ctx, u.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		p, err := token.NewPayLoad(u.Username, u.Role, scopes, 10*time.Minute)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.Set("authorization_payload", p)
		ctx.Next()
	}
}
