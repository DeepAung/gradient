package server

import (
	"net/http"

	"github.com/DeepAung/gradient/cmd/website/internal/config"
	"github.com/DeepAung/gradient/cmd/website/internal/database"
	"github.com/DeepAung/gradient/cmd/website/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func InitServer() error {
	cfg := config.NewConfig()
	db := database.InitPostgresDB(cfg.Database)
	router := chi.NewRouter()
	_ = db

	// s.router.Use(Some Middlewares)

	router.Route("/auths", handlers.NewUserHandler)
	router.Route("/identities", handlers.NewUserHandler)
	router.Route("/persons", handlers.NewUserHandler)
	router.Route("/classrooms", handlers.NewUserHandler)
	router.Route("/tasks", handlers.NewUserHandler)
	router.Route("/submissions", handlers.NewUserHandler)
	router.Route("/testcases", handlers.NewUserHandler)
	router.Route("/tags", handlers.NewUserHandler)
	router.Route("/others", handlers.NewUserHandler)

	return http.ListenAndServe(cfg.App.Address, router)
}
