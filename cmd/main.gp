// package main

// import (
// 	"context"
// 	"log"
// 	"log/slog"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"
// 	"tritchgo/db"
// 	"tritchgo/internal/handlers"
// 	"tritchgo/internal/repository"

// 	"github.com/gofiber/fiber/v2/middleware/compress"
// 	"github.com/gofiber/fiber/v2/middleware/logger"
// 	"github.com/gofiber/fiber/v2/middleware/recover"

// 	"github.com/gofiber/fiber/v2"
// )

// func main() {
// 	LoggerInit()
// 	app := fiber.New(fiber.Config{
// 		Prefork:               true,
// 		CaseSensitive:         true,
// 		StrictRouting:         true,
// 		ReadTimeout:           5 * time.Second,
// 		WriteTimeout:          10 * time.Second,
// 		IdleTimeout:           30 * time.Second,
// 		DisableStartupMessage: true,
// 		EnablePrintRoutes:     false,
// 		StreamRequestBody:     true,
// 	})

// 	app.Use(logger.New())

// 	app.Use(recover.New())
// 	app.Use(compress.New(compress.Config{
// 		Level: compress.LevelBestSpeed,
// 	}))

// 	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
// 	defer stop()

// 	pgdb, err := db.DBConn(ctx)
// 	if err != nil {
// 		log.Fatalf("Fatal conn to db: %v", err)
// 	}
// 	defer pgdb.Close()

// 	rdb := db.RedisDb()
// 	defer rdb.Close()

// 	repo := repository.NewRepository(pgdb)
// 	statsHandle := handlers.NewStatsHandler(repo, pgdb, rdb)

// 	app.Get("/user/stats", statsHandle.GetUserStatsById)
// 	app.Get("/stream/stats", statsHandle.GetStreamStatsById)

// 	go func() {
// 		slog.Info("Server started on :8080")
// 		if err := app.Listen(":8080"); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("Server error: %v", err)
// 		}
// 	}()

// 	<-ctx.Done()
// 	slog.Info("Shutting down server...")

// 	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
// 		slog.Error("Server shutdown error", "error", err)
// 	} else {
// 		slog.Info("Server stopped gracefully")
// 	}
// }
