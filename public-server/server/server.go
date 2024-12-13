package server

import (
	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/modules/auth"
	"github.com/DeepAung/gradient/public-server/modules/tasks"
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/modules/users"
	"github.com/DeepAung/gradient/public-server/modules/views"
	"github.com/DeepAung/gradient/public-server/pkg/storer"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recoverer "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
)

type server struct {
	cfg    *config.Config
	db     *sqlx.DB
	app    *fiber.App
	mid    types.Middleware
	storer storer.Storer
}

func NewServer(
	cfg *config.Config,
	db *sqlx.DB,
	app *fiber.App,
	mid types.Middleware,
	storer storer.Storer,
) *server {
	return &server{
		cfg:    cfg,
		db:     db,
		app:    app,
		mid:    mid,
		storer: storer,
	}
}

func (s *server) Start() {
	s.app.Use(logger.New())
	s.app.Use(recoverer.New())
	s.app.Static("/public", "./public")

	s.setupRoutes()

	s.app.Listen(s.cfg.App.Address)
}

func (s *server) setupRoutes() {
	testGroup := s.app.Group("tests")
	testGroup.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello World")
	})

	apiGroup := s.app.Group("/api")

	usersRepo := users.NewUsersRepo(s.db)
	usersSvc := users.NewUsersSvc(usersRepo, s.storer, s.cfg)
	users.InitUsersHandler(apiGroup, s.mid, usersSvc)

	tasksRepo := tasks.NewTasksRepo(s.db)
	tasksSvc := tasks.NewTasksSvc(tasksRepo)
	tasks.InitTasksHandler(apiGroup, s.mid, tasksSvc)

	authRepo := auth.NewAuthRepo(s.db)
	authSvc := auth.NewAuthSvc(authRepo, usersRepo, s.cfg)
	auth.InitAuthHandler(apiGroup, authSvc, s.cfg)

	views.InitViewsHandler(s.app, s.mid, usersSvc, tasksSvc)
}
