package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/oschwald/geoip2-golang"
	"go.uber.org/zap"

	"github.com/ScMofeoluwa/minalytics/config"
	database "github.com/ScMofeoluwa/minalytics/database/sqlc"
	_ "github.com/ScMofeoluwa/minalytics/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
	config config.Config
	logger *zap.Logger
}

func New(config config.Config, logger *zap.Logger) *Server {
	return &Server{
		router: gin.Default(),
		config: config,
		logger: logger,
	}
}

// @title Minalytics API
// @version 1.0
// @description Analytics API for tracking and managing app data.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func (s *Server) Start() {
	if err := s.migrateDB(); err != nil {
		s.logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	goth.UseProviders(google.New(
		s.config.GoogleClientID,
		s.config.GoogleClientSecret,
		s.config.GoogleClientCallbackUrl,
	), github.New(
		s.config.GithubClientID,
		s.config.GithubClientSecret,
		s.config.GithubClientCallbackUrl,
	))

	geoDB, err := geoip2.Open("database/GeoLite2-City.mmdb")
	if err != nil {
		s.logger.Fatal("Failed to open GeoIP2 database", zap.Error(err))
	}
	defer geoDB.Close()

	connPool, err := pgxpool.New(context.Background(), s.config.DatabaseURL)
	if err != nil {
		s.logger.Fatal("Failed to create connection pool", zap.Error(err))
	}
	defer connPool.Close()

	if err := connPool.Ping(context.Background()); err != nil {
		s.logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	querier := database.New(connPool)
	analyticsService := NewAnalyticsService(querier, geoDB)
	analyticsHandler := NewAnalyticsHandler(analyticsService, s.logger)

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	s.router.GET("/analytics/track", WrapHandler(analyticsHandler.TrackEvent))

	auth := s.router.Group("auth")
	{
		auth.GET(":provider", analyticsHandler.SignIn)
		auth.GET(":provider/callback", WrapHandler(analyticsHandler.Callback))
	}

	apps := s.router.Group("apps")
	apps.Use(JWTMiddleware())
	{
		apps.GET("/", WrapHandler(analyticsHandler.GetApps))
		apps.POST("/", WrapHandler(analyticsHandler.CreateApp))
		apps.PATCH("/:trackingID", WrapHandler(analyticsHandler.UpdateApp))
		apps.DELETE("/:trackingID", WrapHandler(analyticsHandler.UpdateApp))
	}

	analytics := s.router.Group("analytics")
	analytics.Use(JWTMiddleware())
	analytics.Use(AppAccessMiddleware(analyticsService))
	{
		analytics.GET("referrals", WrapHandler(analyticsHandler.GetReferrals))
		analytics.GET("pages", WrapHandler(analyticsHandler.GetPages))
		analytics.GET("browsers", WrapHandler(analyticsHandler.GetBrowsers))
		analytics.GET("countries", WrapHandler(analyticsHandler.GetCountries))
		analytics.GET("devices", WrapHandler(analyticsHandler.GetDevices))
		analytics.GET("os", WrapHandler(analyticsHandler.GetOS))
		analytics.GET("visitors", WrapHandler(analyticsHandler.GetVisitors))
		analytics.GET("pageviews", WrapHandler(analyticsHandler.GetPageViews))
	}

	port := s.config.Port
	s.logger.Info("Starting server", zap.String("port", port))
	s.router.Run(":" + port)
}

func (s *Server) migrateDB() error {
	m, err := migrate.New("file://database/migrations", s.config.DatabaseURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	s.logger.Info("migrations applied successfully")
	return nil
}
