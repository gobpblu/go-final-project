package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

const (
	portKey = "TODO_PORT"
	webDir  = "web"
	indexPath   = "./web/index.html"
)

func main() {
	r := chi.NewRouter()

	port := os.Getenv(portKey)
	if port == "" {
		port = "7540"
	}

	fileServer := http.FileServer(http.Dir(webDir))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, indexPath)
	})
	r.Handle("/*", fileServer)

	fmt.Println("Сервер запущен на http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
