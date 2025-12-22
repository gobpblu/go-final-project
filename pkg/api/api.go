package api

import (
	"go-final-project/pkg/api/nextdate"

	"github.com/go-chi/chi/v5"
)

func Init(r *chi.Mux) {
	r.Get("/api/nextdate", nextdate.NextDateHandler)
}
