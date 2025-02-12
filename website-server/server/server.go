package server

import (
	"github.com/DeepAung/gradient/grader-server/graderconfig"
	"github.com/DeepAung/gradient/grader-server/proto"
	"github.com/DeepAung/gradient/website-server/config"
	"github.com/DeepAung/gradient/website-server/modules/auth"
	"github.com/DeepAung/gradient/website-server/modules/middlewares"
	"github.com/DeepAung/gradient/website-server/modules/submissions"
	"github.com/DeepAung/gradient/website-server/modules/tasks"
	"github.com/DeepAung/gradient/website-server/modules/types"
	"github.com/DeepAung/gradient/website-server/modules/users"
	"github.com/DeepAung/gradient/website-server/modules/views"
	"github.com/DeepAung/gradient/website-server/pkg/storer"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recoverer "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
)

type server struct {
	cfg          *config.Config
	graderCfg    *graderconfig.Config
	db           *sqlx.DB
	app          *fiber.App
	mid          types.Middleware
	storer       storer.Storer
	graderClient proto.GraderClient
}

func NewServer(
	cfg *config.Config,
	graderCfg *graderconfig.Config,
	db *sqlx.DB,
	app *fiber.App,
	storer storer.Storer,
	graderClient proto.GraderClient,
) *server {
	authRepo := auth.NewAuthRepo(db, cfg.App.Timeout)
	usersRepo := users.NewUsersRepo(db, cfg.App.Timeout)
	authSvc := auth.NewAuthSvc(authRepo, usersRepo, cfg)
	mid := middlewares.NewMiddleware(cfg, authSvc)

	return &server{
		cfg:          cfg,
		graderCfg:    graderCfg,
		db:           db,
		app:          app,
		mid:          mid,
		storer:       storer,
		graderClient: graderClient,
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

	usersRepo := users.NewUsersRepo(s.db, s.cfg.App.Timeout)
	usersSvc := users.NewUsersSvc(usersRepo, s.storer, s.cfg)
	users.InitUsersHandler(apiGroup, s.mid, usersSvc)

	authRepo := auth.NewAuthRepo(s.db, s.cfg.App.Timeout)
	authSvc := auth.NewAuthSvc(authRepo, usersRepo, s.cfg)
	auth.InitAuthHandler(apiGroup, authSvc, s.cfg)

	tasksRepo := tasks.NewTasksRepo(s.db, s.cfg.App.Timeout)
	tasksSvc := tasks.NewTasksSvc(tasksRepo)
	tasks.InitTasksHandler(apiGroup, s.mid, tasksSvc)

	submissionsRepo := submissions.NewSubmissionRepo(s.db, s.cfg.App.Timeout)
	submissionsSvc := submissions.NewSubmissionSvc(
		submissionsRepo,
		tasksRepo,
		s.graderClient,
		s.graderCfg,
	)
	submissions.InitSubmissionsHandler(apiGroup, s.mid, submissionsSvc, tasksSvc, s.graderCfg)

	views.InitViewsHandler(s.app, s.mid, usersSvc, tasksSvc, s.graderCfg)
}
