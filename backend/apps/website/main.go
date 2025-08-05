package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/DeepAung/gradient/apps/website/internal/config"
	"github.com/DeepAung/gradient/apps/website/internal/server"
	"github.com/DeepAung/gradient/apps/website/pkg/database"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	api huma.API
	cfg *config.Config
)

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, config *config.Config) {
		initZapLogger(config.App.IsProduction)

		cfg = config

		router := chi.NewMux()
		api = humachi.New(router, huma.DefaultConfig("Gradient", config.App.Version))
		server.InitServer(api, config)

		hooks.OnStart(func() {
			zap.S().Infof("Listening on port %s", config.App.Port)
			if err := http.ListenAndServe(":"+config.App.Port, router); err != nil {
				zap.S().Fatal(err)
			}
		})
	})

	cli.Root().AddCommand(newOpenApiSpecCmd())
	cli.Root().AddCommand(newSetupCmd())

	cli.Run()
}

func initZapLogger(isProduction bool) {
	var config zap.Config
	if isProduction {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Failed to create new zap logger: %v", err)
		return
	}

	zap.ReplaceGlobals(logger)
}

func newOpenApiSpecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI spec",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use downgrade to return OpenAPI 3.0.3 YAML since oapi-codegen doesn't
			// support OpenAPI 3.1 fully yet. Use `.YAML()` instead for 3.1.
			b, err := api.OpenAPI().DowngradeYAML()
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		},
	}
}

//go:embed internal/migrations/*.sql
var embedMigrations embed.FS

func newSetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Setup the gradient application",
		RunE: func(cmd *cobra.Command, args []string) error {
			db := database.InitPostgresDB(&cfg.Database)

			// Migrate the database
			goose.SetBaseFS(embedMigrations)
			if err := goose.SetDialect("postgres"); err != nil {
				return err
			}
			if err := goose.Up(db.DB, "internal/migrations"); err != nil {
				return err
			}

			// TODO:
			// - Create socials (google & github)
			// - Create global gradient with id = `gradient`
			// - Insert tags
			// - Insert problems

			return nil
		},
	}
}
