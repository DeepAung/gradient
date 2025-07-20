package handlers

import "github.com/go-chi/chi/v5"

type userHandler struct{}

func NewUserHandler(r chi.Router) {
	h := &userHandler{}
	_ = h
}
