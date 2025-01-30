package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/oschwald/geoip2-golang"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	database "github.com/ScMofeoluwa/minalytics/database/sqlc"
	_ "github.com/ScMofeoluwa/minalytics/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}
	defer logger.Sync()

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("Error reading config file", zap.Error(err))
	}

	goth.UseProviders(google.New(
		viper.GetString("GOOGLE_CLIENT_ID"),
		viper.GetString("GOOGLE_CLIENT_SECRET"),
		viper.GetString("GOOGLE_CLIENT_CALLBACK_URL"),
	),
	)

	geoDB, err := geoip2.Open("database/GeoLite2-City.mmdb")
	if err != nil {
		logger.Fatal("Failed to open GeoIP2 database", zap.Error(err))
	}
	defer geoDB.Close()

	connPool, err := pgxpool.New(context.Background(), viper.GetString("DATABASE_URL"))
	if err != nil {
		logger.Fatal("Failed to create connection pool", zap.Error(err))
	}
	defer connPool.Close()

	if err := connPool.Ping(context.Background()); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	queries := database.New(connPool)
	analyticsService := NewAnalyticsService(queries, geoDB)
	analyticsHandler := NewAnalyticsHandler(analyticsService, logger)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/analytics/track", WrapHandler(analyticsHandler.TrackEvent))
	r.GET("/account/trackingID", WrapHandler(analyticsHandler.GetTrackingID))

	auth := r.Group("auth")
	{
		auth.GET(":provider", analyticsHandler.SignIn)
		auth.GET(":provider/callback", WrapHandler(analyticsHandler.Callback))
	}

	analytics := r.Group("analytics")
	analytics.Use(JWTMiddleware())
	{
		analytics.GET("referrals", WrapHandler(analyticsHandler.GetReferrals))
		analytics.GET("pages", WrapHandler(analyticsHandler.GetPages))
		analytics.GET("browsers", WrapHandler(analyticsHandler.GetBrowsers))
		analytics.GET("countries", WrapHandler(analyticsHandler.GetCountries))
		analytics.GET("devices", WrapHandler(analyticsHandler.GetDevices))
		analytics.GET("os", WrapHandler(analyticsHandler.GetOS))
	}

	port := viper.GetString("PORT")
	log.Printf("listening on port: %s\n", port)
	r.Run(":" + port)
}
