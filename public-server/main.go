package main

import (
	"flag"

	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/database"
	"github.com/DeepAung/gradient/public-server/modules/middlewares"
	"github.com/DeepAung/gradient/public-server/pkg/graderclient"
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

	// graderClient, conn, err := graderclient.NewGraderClient(
	// 	cfg.App.Address,
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )
	// if err != nil {
	// 	log.Fatal("init grader client error: ", err.Error())
	// }
	// defer conn.Close()
	graderClient := graderclient.NewGraderClientMock(10)

	server := server.NewServer(cfg, db, app, mid, storer, graderClient)
	server.Start()
}
