package main

import (
	"flag"
	"log"

	"github.com/DeepAung/gradient/grader-server/graderconfig"
	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/database"
	"github.com/DeepAung/gradient/public-server/pkg/graderclient"
	"github.com/DeepAung/gradient/public-server/pkg/storer"
	"github.com/DeepAung/gradient/public-server/server"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	envPath  = flag.String("env", ".env.dev", "env path")
	jsonPath = flag.String("json", "../grader-server/.env.dev.json", "grader config json path")
)

func main() {
	flag.Parse()

	graderCfg := graderconfig.NewConfig(*jsonPath)
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

	server := server.NewServer(cfg, graderCfg, db, app, storer, graderClient)
	server.Start()
}
