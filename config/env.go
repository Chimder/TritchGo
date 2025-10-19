package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type EnvVars struct {
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	DBUrl        string `env:"DB_URL"`
	Redis        string `env:"REDIS"`
	Debug        bool   `env:"DEBUG"`
}

func LoadEnv() *EnvVars {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Err load env %v", err)
	}

	cwd, _ := os.Getwd()
	_ = godotenv.Load(filepath.Join(cwd, ".env"))
	_ = godotenv.Load(filepath.Join(cwd, "..", ".env"))
	_ = godotenv.Load(filepath.Join(cwd, "../..", ".env"))

	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..")
	_ = godotenv.Load(filepath.Join(projectRoot, ".env"))

	cfg := &EnvVars{
		ClientID:     getEnv("CLIENT_ID", ""),
		ClientSecret: getEnv("CLIENT_SECRET", ""),
		DBUrl:        getEnv("DB_URL", ""),
		Debug:        getEnv("DEBUG", "") == "true",
	}

	return cfg
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

// package config

// import (
// 	"log/slog"
// 	"os"
// 	"path/filepath"

// 	"github.com/joho/godotenv"
// )

// type EnvVars struct {
// 	ClientID     string `env:"CLIENT_ID"`
// 	ClientSecret string `env:"CLIENT_SECRET"`
// 	DBUrl        string `env:"DB_URL"`
// 	Debug        bool   `env:"DEBUG"`
// }

// func LoadEnv() *EnvVars {
// 	cwd, _ := os.Getwd()

// 	_ = godotenv.Load(filepath.Join(cwd, ".env"))

// 	_ = godotenv.Load(filepath.Join(cwd, "..", ".env"))
// 	if err := godotenv.Load(); err != nil {
// 		slog.Error("ENV not found")
// 	}

// 	cfg := &EnvVars{
// 		ClientID:     getEnv("CLIENT_ID", ""),
// 		ClientSecret: getEnv("CLIENT_SECRET", ""),
// 		DBUrl:        getEnv("DB_URL", ""),
// 		Debug:        getEnv("DEBUG", "") == "true",
// 	}

// 	return cfg
// }

// func getEnv(key, defaultVal string) string {
// 	if value := os.Getenv(key); value != "" {
// 		return value
// 	}

// 	return defaultVal
// }
