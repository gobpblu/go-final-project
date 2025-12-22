package server

import (
	"go-final-project/pkg/api"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	webDir    = "web"
	indexPath = "./web/index.html"
)

func Run(r *chi.Mux) {
	fileServer := http.FileServer(http.Dir(webDir))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, indexPath)
	})
	r.Handle("/*", fileServer)

	api.Init(r)
}
