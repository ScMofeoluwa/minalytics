package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			ctx.Abort()
			return
		}

		claims, err := VerifyJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		userIDStr, ok := claims["sub"].(string)
		if !ok || userIDStr == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid UserID in token"})
			ctx.Abort()
			return
		}

		userID, _ := uuid.Parse(userIDStr)
		ctx.Set("userID", userID)
		ctx.Next()
	}
}

func WrapHandler(handler func(*gin.Context) APIResponse) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response := handler(ctx)
		ctx.JSON(response.statusCode, response)
	}
}
