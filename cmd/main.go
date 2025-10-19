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
	"tritchgo/config"
	"tritchgo/db"
	"tritchgo/internal/handlers"
	"tritchgo/internal/kafka"
	"tritchgo/internal/nats"
	"tritchgo/internal/repository"
	"tritchgo/internal/routers"
)

func main() {
	env := config.LoadEnv()
	LoggerInit(env)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	twitchHandle := handlers.NewTwitchHandle(env)
	pgdb, err := db.DBConn(ctx, env.DBUrl)
	if err != nil {
		log.Fatalf("Fatal conn to db: %v", err)
	}
	defer pgdb.Close()

	repo := repository.NewRepository(pgdb)

	els := db.NewElasticDB()
	err = els.Mapping()
	if err != nil {
		log.Fatalf("Fatal conn to elasticsearch: %v", err)
	}

	natsStream := nats.NewNatsProducer(ctx)
	defer natsStream.Close()

	rdb := db.RedisDb(env.Redis)
	defer rdb.Close()

	kafkaProducer := kafka.NewKafkaWriters()
	defer kafkaProducer.Close()

	go StartGRPCServer(pgdb)
	go NewTwitchScheduler(ctx, pgdb, els.GetClient(), natsStream, repo).StartFetchLoop(twitchHandle)

	r := routers.NewRouter(repo, pgdb, rdb, kafkaProducer, els.GetClient())

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	} else {
		slog.Info("Server stopped gracefully")
	}
}
