package main

import (
	"fmt"
	"go-final-project/pkg/constants"
	"go-final-project/pkg/db"
	"go-final-project/pkg/server"
	"go-final-project/pkg/utils"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные из .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	r := chi.NewRouter()

	port := utils.GetEnvOrDefault(constants.PortKey, constants.DefaultPort)
	dbFile := utils.GetEnvOrDefault(constants.DbKey, constants.DefaultDbName)

	err = db.Init(dbFile)
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
