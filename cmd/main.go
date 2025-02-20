package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"tritchgo/db"
	"tritchgo/internal/handlers"
	"tritchgo/internal/routers"
)

func main() {
	context := context.Background()

	twitchHandle := handlers.NewTwitchHandle()
	db, err := db.DBConn(context)
	if err != nil {
		log.Fatalf("Fatal conn to db: %v", err)
	}

	go StartGRPCServer(db)
	go NewTwitchSheduler(context, db).StartFetchLoop(twitchHandle)

	r := routers.NewRouter(db)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server started on :8080")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
