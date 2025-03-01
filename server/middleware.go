package server

import (
	"net/http"
	"strings"

	types "github.com/ScMofeoluwa/minalytics/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

		userID := uuid.MustParse(userIDStr)
		ctx.Set("userID", userID)
		ctx.Next()
	}
}

func AppAccessMiddleware(s types.AnalyticsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
			ctx.Abort()
			return
		}
		user := userID.(uuid.UUID)

		trackingID_ := ctx.Query("trackingID")
		if trackingID_ == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "trackingID is required"})
			ctx.Abort()
			return
		}

		trackingID := uuid.MustParse(trackingID_)

		err := s.ValidateAppAccess(ctx, user, trackingID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("trackingID", trackingID)
		ctx.Next()
	}
}

func WrapHandler(handler func(*gin.Context) types.APIResponse) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response := handler(ctx)
		ctx.JSON(response.StatusCode, response)
	}
}
