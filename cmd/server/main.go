package main

import (
	"log"

	"github.com/ScMofeoluwa/minalytics/config"
	"github.com/ScMofeoluwa/minalytics/server"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Config error: %s", zap.Error(err))
	}

	server := server.New(cfg, logger)
	server.Start()
}
