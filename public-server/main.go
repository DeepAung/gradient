package main

import (
	"flag"

	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/database"
	"github.com/DeepAung/gradient/public-server/modules/middlewares"
	"github.com/DeepAung/gradient/public-server/pkg/storer"
	"github.com/DeepAung/gradient/public-server/server"
	"github.com/gofiber/fiber/v2"
)

var envPath = flag.String("env", "", "env file")

func main() {
	flag.Parse()

	cfg := config.NewConfig(*envPath)
	db := database.InitDB(cfg.App.DbUrl)
	app := fiber.New()

	mid := middlewares.NewMiddleware(cfg)
	storer := storer.NewGcpStorer(cfg)

	server := server.NewServer(cfg, db, app, mid, storer)
	server.Start()
}
