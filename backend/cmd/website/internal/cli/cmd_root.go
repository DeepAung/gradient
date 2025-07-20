package cli

import (
	"os"

	"github.com/DeepAung/gradient/cmd/website/internal/config"
	"github.com/DeepAung/gradient/cmd/website/internal/database"
	"github.com/DeepAung/gradient/cmd/website/internal/migrations"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
)

func Execute() error {
	cfg := config.NewConfig()
	db := database.InitPostgresDB(cfg.Database)
	app := &cli.App{
		Name: "gradient-website",
		Commands: []*cli.Command{
			newDBCommand(migrate.NewMigrator(db, migrations.Migrations)),
			newServeCommand(),
		},
	}
	return app.Run(os.Args)
}
