package main

import (
	"flag"

	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/database"
	"github.com/DeepAung/gradient/public-server/server"
	"github.com/gofiber/fiber/v2"
)

var envPath = flag.String("env", "", "env file")

func main() {
	flag.Parse()

	cfg := config.NewConfig(*envPath)
	cfg.Print()
	db := database.InitDB(cfg.App.DbUrl)
	app := fiber.New()

	server := server.NewServer(cfg, db, app)
	server.Start()
}
