package api

import (
	"go-final-project/pkg/api/middlewares"
	"go-final-project/pkg/api/nextdate"
	"go-final-project/pkg/api/signin"
	"go-final-project/pkg/api/tasks"

	"github.com/go-chi/chi/v5"
)

func Init(r *chi.Mux) {
	r.Post("/api/signin", signin.SignInHandler)
	r.Get("/api/nextdate", nextdate.NextDateHandler)

	r.Route("/api/task", func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Post("/", tasks.AddTaskHandler)
		r.Get("/", tasks.GetTaskHandler)
		r.Put("/", tasks.UpdateTaskHandler)
		r.Delete("/", tasks.DeleteTaskHandler)
		r.Post("/done", tasks.TaskDoneHandler)
	})

	r.Route("/api/tasks", func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Get("/", tasks.GetTasksHandler)
	})
}
