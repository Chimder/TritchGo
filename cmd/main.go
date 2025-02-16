package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"tritchgo/db"
	"tritchgo/internal/routers"
)

type twitchToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func main() {
	context := context.Background()

	// twitchHandle := handlers.NewTwitchHandle()

	db, err := db.DBConn(context)
	if err != nil {
		log.Fatalf("Fatal conn to db: %v", err)
	}
	// go func() {
	// 	for {
	// 		stats := db.Stat()
	// 		log.Printf("Pool stats: TotalConns=%d, IdleConns=%d, AcquiredConns=%d", stats.TotalConns(), stats.IdleConns(), stats.AcquiredConns())
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	// go StartFetchLoop(context, twitchHandle, db)

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
