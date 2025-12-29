package api

import (
	"go-final-project/pkg/api/nextdate"
	"go-final-project/pkg/api/tasks"

	"github.com/go-chi/chi/v5"
)

func Init(r *chi.Mux) {
	r.Get("/api/nextdate", nextdate.NextDateHandler)
	r.Post("/api/task", tasks.AddTaskHandler)
	r.Get("/api/tasks", tasks.GetTasksHandler)
}
