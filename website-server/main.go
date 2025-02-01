package main

import (
	"flag"
	"log"

	_ "embed"

	"github.com/DeepAung/gradient/grader-server/graderconfig"
	"github.com/DeepAung/gradient/website-server/config"
	"github.com/DeepAung/gradient/website-server/database"
	"github.com/DeepAung/gradient/website-server/pkg/graderclient"
	"github.com/DeepAung/gradient/website-server/pkg/storer"
	"github.com/DeepAung/gradient/website-server/server"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:embed graderconfig.json
var graderConfigFile []byte

var envPath = flag.String("env", ".env.dev", "env path")

func main() {
	flag.Parse()

	graderCfg := graderconfig.NewConfig(graderConfigFile)
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
