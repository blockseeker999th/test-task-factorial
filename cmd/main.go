package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blockseeker999th/test-task-factorial/pkg/config"
	"github.com/blockseeker999th/test-task-factorial/pkg/db"
	"github.com/blockseeker999th/test-task-factorial/pkg/db/storage"
	savefactorial "github.com/blockseeker999th/test-task-factorial/pkg/handlers/saveFactorial"
	validateinput "github.com/blockseeker999th/test-task-factorial/pkg/middleware/validateInput"
	slogerr "github.com/blockseeker999th/test-task-factorial/pkg/utils/logger/slogErr"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	err := godotenv.Load("./../.env")
	if err != nil {
		fmt.Println("Error loading .env variables", err)
	}

	cfg := config.MustLoad()

	log := setupLogger(os.Getenv("ENV"))
	log.Info("starting factorial app", slog.String("env", os.Getenv("ENV")))

	db, err := db.ConnectDB(cfg).InitNewPostgreSQLStorage()
	st := storage.NewStorage(db)

	if err != nil {
		log.Error("Failed to init PostgreSQL Database", slogerr.Err(err))
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error("Error closing DB: ", slogerr.Err(err))
		}
	}()

	router := httprouter.New()

	router.POST("/calculate", validateinput.ValidateInput(log, savefactorial.New(log, st)))

	server := http.Server{
		Addr:         cfg.HTTPServer.Address,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		Handler:      router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))
		if err := server.ListenAndServe(); err != nil {
			log.Error("server stopped", err)
		}
	}()

	<-quit

	log.Warn("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", err)
	}

	log.Info("Server gracefully stopped")
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
