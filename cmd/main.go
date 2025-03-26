package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
	"tritchgo/db"
	"tritchgo/internal/handlers"
	"tritchgo/internal/routers"
)

func main() {
	context := context.Background()
	opts := &slog.HandlerOptions{
		// AddSource: true,
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	slog.SetDefault(logger)

	twitchHandle := handlers.NewTwitchHandle()
	pgdb, err := db.DBConn(context)
	if err != nil {
		log.Fatalf("Fatal conn to db: %v", err)
	}
	rdb := db.RedisDb()

	go StartGRPCServer(pgdb)
	go NewTwitchSheduler(context, pgdb).StartFetchLoop(twitchHandle)

	r := routers.NewRouter(pgdb, rdb)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.Info("Server started on :8080")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
