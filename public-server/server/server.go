package server

import (
	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/modules/auth"
	"github.com/DeepAung/gradient/public-server/modules/middlewares"
	"github.com/DeepAung/gradient/public-server/modules/tasks"
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/modules/users"
	"github.com/DeepAung/gradient/public-server/modules/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recoverer "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
)

type server struct {
	cfg *config.Config
	db  *sqlx.DB
	app *fiber.App
	mid types.Middleware
}

func NewServer(cfg *config.Config, db *sqlx.DB, app *fiber.App) *server {
	return &server{
		cfg: cfg,
		db:  db,
		app: app,
		mid: middlewares.NewMiddleware(cfg),
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
	usersSvc := users.NewUsersSvc(usersRepo)

	tasksGroup := apiGroup.Group("/tasks")
	tasksRepo := tasks.NewTasksRepo(s.db)
	tasksSvc := tasks.NewTasksSvc(tasksRepo)
	tasks.InitTasksHandler(tasksGroup, s.mid, tasksSvc)

	authGroup := apiGroup.Group("/auth")
	authRepo := auth.NewAuthRepo(s.db)
	authSvc := auth.NewAuthSvc(authRepo, usersRepo, s.cfg)
	auth.InitAuthHandler(authGroup, authSvc, s.cfg)

	views.InitViewsHandler(s.app, s.mid, usersSvc, tasksSvc)
}
