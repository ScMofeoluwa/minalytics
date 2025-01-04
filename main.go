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

	database "github.com/ScMofeoluwa/minalytics/database/sqlc"
)

func main() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	goth.UseProviders(google.New(
		viper.GetString("GOOGLE_CLIENT_ID"),
		viper.GetString("GOOGLE_CLIENT_SECRET"),
		viper.GetString("GOOGLE_CLIENT_CALLBACK_URL"),
	),
	)

	geoDB, err := geoip2.Open("database/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatalf("Failed to open GeoIP2 database: %v", err)
	}
	defer geoDB.Close()

	connPool, err := pgxpool.New(context.Background(), viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer connPool.Close()

	queries := database.New(connPool)
	analyticsService := NewAnalyticsService(queries, geoDB)
	analyticsHandler := NewAnalyticsHandler(analyticsService)

	r := gin.Default()
	r.GET("/", analyticsHandler.Home)
	r.GET("/auth/:provider", analyticsHandler.SignIn)
	r.GET("/auth/:provider/callback", analyticsHandler.Callback)
	r.GET("/track", analyticsHandler.TrackEvent)

	port := viper.GetString("PORT")
	log.Printf("listening on port: %s\n", port)
	r.Run(":" + port)
}
