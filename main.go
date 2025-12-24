package main

import (
	"fmt"
	"go-final-project/pkg/db"
	"go-final-project/pkg/server"
	"go-final-project/pkg/utils"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	portKey       = "TODO_PORT"
	dbKey         = "TODO_DBFILE"
	defaultPort   = "7540"
	defaultDbName = "scheduler.db"
)

func main() {
	r := chi.NewRouter()

	port := utils.GetEnvOrDefault(portKey, defaultPort)
	dbFile := utils.GetEnvOrDefault(dbKey, defaultDbName)

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
