package server

import (
	"github.com/DeepAung/gradient/apps/website/internal/config"
	"github.com/danielgtaylor/huma/v2"
)

func InitServer(api huma.API, cfg *config.Config) error {
	// cfg := config.NewConfig()
	// db := db.InitPostgresDB(&cfg.Database)
	// router := chi.NewRouter()
	// _ = db
	//
	// // s.router.Use(Some Middlewares)
	//
	// router.Route("/auths", handler.NewUserHandler)
	// router.Route("/identities", handler.NewUserHandler)
	// router.Route("/persons", handler.NewUserHandler)
	// router.Route("/classrooms", handler.NewUserHandler)
	// router.Route("/tasks", handler.NewUserHandler)
	// router.Route("/submissions", handler.NewUserHandler)
	// router.Route("/testcases", handler.NewUserHandler)
	// router.Route("/tags", handler.NewUserHandler)
	// router.Route("/others", handler.NewUserHandler)
	//
	// return http.ListenAndServe(cfg.App.Address, router)

	return nil
}
