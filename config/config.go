package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config struct holds all the configuration settings for your application.
type Config struct {
	AppName string
	AppEnv  string
	AppPort string
	DBHost  string
	DBPort  string
	DBUser  string
	DBPass  string
	DBName  string
	Buckets map[string]BucketConfig
}

// BucketConfig holds the configuration for an individual bucket.
type BucketConfig struct {
	AccessKey  string
	SecretKey  string
	BucketName string
}

// LoadConfig loads the configuration from environment variables and config files.
func LoadConfig() (*Config, error) {
	// Load environment variables from a .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	viper.AutomaticEnv() // Automatically override values with environment variables

	// Initialize the buckets map
	buckets := make(map[string]BucketConfig)
	buckets["bucket1"] = BucketConfig{
		AccessKey:  getEnv("BUCKET1_ACCESS_KEY", ""),
		SecretKey:  getEnv("BUCKET1_SECRET_KEY", ""),
		BucketName: getEnv("BUCKET1_BUCKET_NAME", ""),
	}
	buckets["bucket2"] = BucketConfig{
		AccessKey:  getEnv("BUCKET2_ACCESS_KEY", ""),
		SecretKey:  getEnv("BUCKET2_SECRET_KEY", ""),
		BucketName: getEnv("BUCKET2_BUCKET_NAME", ""),
	}

	config := &Config{
		AppName: getEnv("APP_NAME", "SecureStore"),
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),
		DBHost:  getEnv("DB_HOST", "localhost"),
		DBPort:  getEnv("DB_PORT", "5432"),
		DBUser:  getEnv("DB_USER", "user"),
		DBPass:  getEnv("DB_PASS", "password"),
		DBName:  getEnv("DB_NAME", "securedocs"),
		Buckets: buckets,
	}

	return config, nil
}

// getEnv retrieves the value of the environment variable named by the key, or returns the fallback value if the environment variable is not present.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
