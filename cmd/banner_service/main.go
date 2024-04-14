package main

import (
	"context"
	"fmt"
	"github.com/fishmanDK/avito_test_task/internal/clients/rabbitmq"
	"github.com/fishmanDK/avito_test_task/internal/config"
	"github.com/fishmanDK/avito_test_task/internal/handlers"
	"github.com/fishmanDK/avito_test_task/internal/service"
	"github.com/fishmanDK/avito_test_task/internal/storage"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(envLocal)
	logger.Info("setup logger", cfg)
	db, err := storage.MustStorage(cfg.Postgres)
	if err != nil {
		logger.Error("error setup storage", err)
		panic(err)
	}

	srvс, err := service.NewService(logger, cfg.Clients, db)
	if err != nil {
		logger.Error("error setup service", err)
		panic(err)
	}

	handl := handlers.MustHandlers(srvс, logger)

	routs := handl.InitRouts()

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      routs,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	rq, err := rabbitmq.NewRabbitMQConsumer(srvс.DeleteService) // Замените на адрес вашего сервера NATS
	if err != nil {
		log.Fatalf("Failed to create NUTS: %v", err)
	}
	logger.Info("start RabbitMQ")
	go func() {
		if err := rq.SubscribeAndReadMessage(); err != nil {
			log.Printf("Error subscribing and reading messages: %v", err)
		}
	}()

	logger.Info("banner_service started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	sig := <-stop
	fmt.Printf("Received signal: %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", slog.String("err", err.Error()))
		os.Exit(1)
	}

	logger.Info("Server gracefully stopped")

}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		slogHandler := slog.NewTextHandler(os.Stdout, opts)
		logger = slog.New(slogHandler)
	case envDev:
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		slogHandler := slog.NewJSONHandler(os.Stdout, opts)
		logger = slog.New(slogHandler)
	}

	return logger
}
