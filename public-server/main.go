package main

import (
	"flag"
	"log"

	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/database"
	"github.com/DeepAung/gradient/public-server/pkg/graderclient"
	"github.com/DeepAung/gradient/public-server/pkg/storer"
	"github.com/DeepAung/gradient/public-server/server"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var envPath = flag.String("env", "", "env file")

func main() {
	flag.Parse()

	cfg := config.NewConfig(*envPath)
	db := database.InitDB(cfg.App.DbUrl)
	app := fiber.New()

	storer := storer.NewGcpStorer(cfg)

	graderClient, conn, err := graderclient.NewGraderClient(
		cfg.App.GraderAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("graderclient.NewGraderClient: ", err.Error())
	}
	defer conn.Close()

	server := server.NewServer(cfg, db, app, storer, graderClient)
	server.Start()
}
