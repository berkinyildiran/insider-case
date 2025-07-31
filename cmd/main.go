package main

import (
	"context"
	"github.com/berkinyildiran/insider-case/internal/cache/redis"
	"github.com/berkinyildiran/insider-case/internal/config"
	"github.com/berkinyildiran/insider-case/internal/database"
	"github.com/berkinyildiran/insider-case/internal/message"
	"github.com/berkinyildiran/insider-case/internal/scheduler"
	"github.com/berkinyildiran/insider-case/internal/sender"
	"github.com/berkinyildiran/insider-case/internal/server"
	"github.com/berkinyildiran/insider-case/internal/transporter/http"
	"github.com/berkinyildiran/insider-case/internal/validator"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Transporter and Validator
	newTransporter := http.NewHttp()
	newValidator := validator.NewValidator()

	// Config
	newConfig := config.NewConfig(newValidator)

	if err := newConfig.Load("config", "./internal/config"); err != nil {
		log.Fatal(err)
	}
	if err := newConfig.Validate(); err != nil {
		log.Fatal(err)
	}

	// Cache and Database
	newCache := redis.NewRedis(newConfig.Cache, ctx)
	newDatabase := database.NewDatabase(newConfig.Database, ctx)

	if err := newDatabase.Connect(); err != nil {
		log.Fatal(err)
	}

	// Database Migration
	models := []any{&message.Message{}}
	if err := newDatabase.Migrate(models...); err != nil {
		log.Fatal(err)
	}

	// Repository
	newRepository := message.NewRepository(newDatabase, ctx)

	// Sender and Scheduler
	newSender := sender.NewSender(newCache, newConfig.Sender, newRepository, newTransporter)
	newScheduler := scheduler.NewScheduler(newConfig.Scheduler, newSender.Run)

	// Handler
	newHandler := message.NewHandler(newRepository, newScheduler, newValidator)

	// Router
	router := server.NewRouter(newConfig.Server, newHandler, newValidator, ctx)
	router.Setup()

	// Start server
	go func() {
		if err := router.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal (e.g., CTRL+C or SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Shutdown signal received. Initiating graceful shutdown...")

	// === 1. Stop HTTP Server ===
	log.Println("[INFO] Stopping HTTP server...")
	if err := router.Stop(); err != nil {
		log.Printf("[ERROR] Failed to stop HTTP server: %v", err)
	} else {
		log.Println("[INFO] HTTP server stopped")
	}

	// === 2. Stop Background Scheduler ===
	log.Println("[INFO] Stopping scheduler...")
	newScheduler.Stop()

	// === 3. Close Cache (Redis) ===
	log.Println("[INFO] Closing Cache connection...")
	if err := newCache.Close(); err != nil {
		log.Printf("[ERROR] Failed to close Cache: %v", err)
	} else {
		log.Println("[INFO] Cache connection closed")
	}

	// === 4. Close Database Connection ===
	log.Println("[INFO] Closing database connection...")
	if err := newDatabase.Close(); err != nil {
		log.Printf("[ERROR] Failed to close database: %v", err)
	} else {
		log.Println("[INFO] Database connection closed")
	}

	log.Println("[INFO] Graceful shutdown completed successfully")
}
