package cli

import (
	"github.com/DeepAung/gradient/cmd/website/internal/server"
	"github.com/urfave/cli/v2"
)

func newServeCommand() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "serve server",
		Action: func(ctx *cli.Context) error {
			return server.InitServer()
		},
	}
}
