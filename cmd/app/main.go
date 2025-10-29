package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vijitt-server/internal/db"
	"vijitt-server/internal/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		os.Exit(0)
	}
	err = db.Init()
	if err != nil {
		log.Println("Server startup failed! Database connection error.")
		os.Exit(0)
	}

	e := echo.New()

	if os.Getenv("LEVEL") != "DEV" {
		// Open a file for logging
		logFile, err := os.OpenFile("/app/logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
		log.SetFlags(log.LstdFlags | log.Lshortfile)

		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "Time=${time_rfc3339}, Method=${method}, Uri=${uri}, Status=${status}\n",
			Output: logFile,
		}))
	} else {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "Time=${time_rfc3339}, Method=${method}, Uri=${uri}, Status=${status}\n",
		}))
	}

	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	routes.HealthRoutes(e)
	// admin_api := e.Group("/admin")

	// routes.AdminRoutes(admin_api)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
