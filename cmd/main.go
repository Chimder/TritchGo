package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"tritchgo/internal/handlers"
	"tritchgo/internal/queries"
	"tritchgo/internal/routers"
	"tritchgo/sqlc"
)

type twitchToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func main() {
	context := context.Background()

	twitchHandle := handlers.NewTwitchHandle()

	db, err := sqlc.DBConn(context)
	sqlc := queries.New(db)
	if err != nil {
		log.Fatalf("Fatal conn to db: %v", err)
	}

	go StartFetchLoop(context, twitchHandle, sqlc)

	r := routers.NewRouter(sqlc)
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
