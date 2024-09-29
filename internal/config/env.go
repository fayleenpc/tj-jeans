package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
	JWTRefresh             string
	SMTP_User              string
	SMTP_Password          string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	// init_vault()
	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", ""),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:                 getEnv("DB_NAME", "tj-jeans"),
		JWTSecret:              getEnv("JWT_SECRET", ""),
		JWTRefresh:             getEnv("JWT_REFRESH", ""),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 300),
		SMTP_User:              getEnv("SMTP_USER", ""),
		SMTP_Password:          getEnv("SMTP_PASSWORD", ""),
	}
}

func init_vault() {
	// Initialize Vault client
	config := api.DefaultConfig()
	config.Address = "http://127.0.0.1:8200" // Your Vault server address

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("unable to create Vault client: %v", err)
	}

	// Set the Vault token (use the root token from dev mode)
	client.SetToken("hvs.qHFXsWObJeaPCrxC52Q19rLS") // Use the root token provided by `vault server -dev`

	// Path and secret for storing environment variables
	path := "secret/data/tj-jeans/config"

	// Read secret from Vault
	secretData, err := client.Logical().Read(path)
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}

	if secretData == nil {
		log.Fatalf("secret not found")
	}

	data, ok := secretData.Data["data"].(map[string]interface{})
	if !ok {
		log.Fatalf("invalid secret data format")
	}

	// Set environment variables
	for key, value := range data {
		os.Setenv(key, fmt.Sprintf("%v", value))
		fmt.Printf("Set environment variable: %s=%s\n", key, value)
	}

	// Verify environment variables
	fmt.Println("\nEnvironment Variables:")
	fmt.Printf("DB_USER: %s\n", os.Getenv("DB_USER"))
	fmt.Printf("DB_PASSWORD: %s\n", os.Getenv("DB_PASSWORD"))
	fmt.Printf("DB_HOST: %s\n", os.Getenv("DB_HOST"))
	fmt.Printf("DB_PORT: %s\n", os.Getenv("DB_PORT"))
	fmt.Printf("DB_NAME: %s\n", os.Getenv("DB_NAME"))
	fmt.Printf("JWT_SECRET: %s\n", os.Getenv("JWT_SECRET"))
	fmt.Printf("JWT_REFRESH: %s\n", os.Getenv("JWT_REFRESH"))
	fmt.Printf("SMTP_USER: %s\n", os.Getenv("SMTP_USER"))
	fmt.Printf("SMTP_PASSWORD: %s\n", os.Getenv("SMTP_PASSWORD"))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
