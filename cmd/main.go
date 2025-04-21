package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tritchgo/db"
	"tritchgo/internal/handlers"
	"tritchgo/internal/kafka"
	"tritchgo/internal/routers"
)

func main() {
	LoggerInit()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	twitchHandle := handlers.NewTwitchHandle()
	pgdb, err := db.DBConn(ctx)
	if err != nil {
		log.Fatalf("Fatal conn to db: %v", err)
	}
	defer pgdb.Close()

	rdb := db.RedisDb()
	defer rdb.Close()

	kafkaProducer := kafka.NewProdKafka()
	defer kafkaProducer.Close()

	go StartGRPCServer(pgdb)
	go NewTwitchSheduler(ctx, pgdb).StartFetchLoop(twitchHandle)

	r := routers.NewRouter(pgdb, rdb, kafkaProducer.GetWriter())
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	} else {
		slog.Info("Server stopped gracefully")
	}
}
