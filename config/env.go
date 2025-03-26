package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type EnvVars struct {
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	DBUrl        string `env:"DB_URL"`
}

func LoadEnv() *EnvVars {
	_ = godotenv.Load()
	cfg := EnvVars{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}

// import (
// 	"log"
// 	"os"

// 	"github.com/joho/godotenv"
// )

// type EnvVars struct {
// 	CLIENT_ID     string
// 	CLIENT_SECRET string
// 	DB_URL        string
// }

// func LoadEnv() EnvVars {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Println("Error loading .env file")
// 	}

// 	client_id := os.Getenv("CLIENT_ID")
// 	client_secret := os.Getenv("CLIENT_SECRET")
// 	db_url := os.Getenv("DB_URL")

// 	// isProdStr := os.Getenv("IS_PROD")
// 	// isProd, err := strconv.ParseBool(isProdStr)
// 	// if err != nil {
// 	// 	isProd = false
// 	// }

// 	return EnvVars{
// 		CLIENT_ID:     client_id,
// 		CLIENT_SECRET: client_secret,
// 		DB_URL:        db_url,
// 	}
// }
