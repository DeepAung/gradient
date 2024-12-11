package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App *AppConfig
	Jwt *JwtConfig
}

type AppConfig struct {
	Address       string
	Timeout       time.Duration
	BodyLimit     string
	DbUrl         string
	GcpBucketName string
}

type JwtConfig struct {
	SecretKey      []byte
	AccessExpires  time.Duration
	RefreshExpires time.Duration
}

func NewConfig(envPath string) *Config {
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("godotenv.Load: ", err)
	}

	return &Config{
		App: &AppConfig{
			Address:       os.Getenv("ADDRESS"),
			Timeout:       getenvDuration("TIMEOUT"),
			BodyLimit:     os.Getenv("BODY_LIMIT"),
			DbUrl:         os.Getenv("DB_URL"),
			GcpBucketName: os.Getenv("GCP_BUCKET_NAME"),
		},
		Jwt: &JwtConfig{
			SecretKey:      []byte(os.Getenv("SECRET_KEY")),
			AccessExpires:  getenvDuration("ACCESS_EXPIRES"),
			RefreshExpires: getenvDuration("REFRESH_EXPIRES"),
		},
	}
}

func (c *Config) Print() {
	fmt.Printf("================ Config ================\n")
	fmt.Printf("%+v\n", c.App)
	fmt.Printf("%+v\n", c.Jwt)
	fmt.Printf("========================================\n")
}

func getenvDuration(key string) time.Duration {
	s := os.Getenv(key)
	if s == "" {
		return 0
	}

	num, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("config.go: convert string to int error. (%q=%q)\n", s, num)
	}

	return time.Duration(num) * time.Second
}
