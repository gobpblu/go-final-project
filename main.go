package main

import (
	"fmt"
	"go-final-project/pkg/db"
	"go-final-project/pkg/server"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

const (
	portKey = "TODO_PORT"
	dbKey   = "TODO_DBFILE"
)

func main() {
	r := chi.NewRouter()

	port := os.Getenv(portKey)
	if port == "" {
		port = "7540"
	}

	dbFile := os.Getenv(dbKey)
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка при открытии БД: %s", err.Error())
		return
	}

	server.Run(r)

	fmt.Println("Сервер запущен на http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
