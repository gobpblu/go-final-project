package main

import (
	"log"
	"net/http"

	"go-final-project/pkg/constants"
	"go-final-project/pkg/db"
	"go-final-project/pkg/server"
	"go-final-project/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные из .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки .env файла")
	}

	r := chi.NewRouter()

	port := utils.GetEnvOrDefault(constants.PortKey, constants.DefaultPort)
	dbFile := utils.GetEnvOrDefault(constants.DbKey, constants.DefaultDbName)
	password := utils.GetEnvOrDefault(constants.PasswordKey, "")

	err, database := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка при открытии БД: %s", err.Error())
		return
	}
	defer database.Close()

	server.Run(r, password)

	log.Println("Сервер запущен на http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
