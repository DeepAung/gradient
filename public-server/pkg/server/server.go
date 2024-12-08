package server

import (
	"github.com/DeepAung/gradient/public-server/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recoverer "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
)

type server struct {
	cfg *config.Config
	db  *sqlx.DB
	app *fiber.App
}

func NewServer(cfg *config.Config, db *sqlx.DB, app *fiber.App) *server {
	return &server{
		cfg: cfg,
		db:  db,
		app: app,
	}
}

func (s *server) Start() {
	s.app.Use(logger.New())
	s.app.Use(recoverer.New())

	// s.app.Use(jwtware.New(jwtware.Config{
	// 	SigningKey:  jwtware.SigningKey{Key: s.cfg.Jwt.SecretKey},
	// 	TokenLookup: "cookie:access-token",
	// }))

	// s.App.Static("/static", "./static")

	s.setupRoutes()

	s.app.Listen(s.cfg.App.Address)
}

func (s *server) setupRoutes() {
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello World")
	})
}
