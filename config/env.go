package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVars struct {
	CLIENT_ID     string
	CLIENT_SECRET string
	DB_URL        string
}

func LoadEnv() EnvVars {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	client_id := os.Getenv("CLIENT_ID")
	client_secret := os.Getenv("CLIENT_SECRET")
	db_url := os.Getenv("DB_URL")

	// isProdStr := os.Getenv("IS_PROD")
	// isProd, err := strconv.ParseBool(isProdStr)
	// if err != nil {
	// 	isProd = false
	// }

	return EnvVars{
		CLIENT_ID:     client_id,
		CLIENT_SECRET: client_secret,
		DB_URL:        db_url,
	}
}

// package config

// import (
// 	"log"
// 	"os"
// )

// type EnvVars struct {
// 	CLIENT_ID     string
// 	CLIENT_SECRET string
// 	DB_URL        string
// }

// func LoadEnv() EnvVars {
// 	clientID := os.Getenv("CLIENT_ID")
// 	if clientID == "" {
// 		log.Print("CLIENT_ID is not set")
// 	}

// 	clientSecret := os.Getenv("CLIENT_SECRET")
// 	if clientSecret == "" {
// 		log.Print("CLIENT_SECRET is not set")
// 	}

// 	dbURL := os.Getenv("DB_URL")
// 	if dbURL == "" {
// 		log.Print("DB_URL is not set")
// 	}
// 	log.Printf("db", dbURL)
// 	log.Printf("id", clientID)
// 	log.Printf("cliek", clientSecret)

// 	return EnvVars{
// 		CLIENT_ID:     clientID,
// 		CLIENT_SECRET: clientSecret,
// 		DB_URL:        dbURL,
// 	}
// }
