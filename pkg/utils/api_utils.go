package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type JsonData map[string]any

func WriteJSON(w http.ResponseWriter, status int, data JsonData, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func WriteBadRequestError(w http.ResponseWriter) {
	writeErrorResponse(w, http.StatusBadRequest, "неправильно указано тело запроса")
}

func WriteFailedValidationError(w http.ResponseWriter, err error) {
	writeErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
}

func WriteInternalServerError(w http.ResponseWriter) {
	writeErrorResponse(w, http.StatusInternalServerError, "что-то пошло не так")
}

func writeErrorResponse(w http.ResponseWriter, status int, message any) {
	env := JsonData{"error": message}

	err := WriteJSON(w, status, env, nil)
	if err != nil {
		log.Fatal(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
