package signin

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"go-final-project/pkg/utils"
)

const (
	passwordKey = "TODO_PASSWORD"
)

var (
	IncorrectPasswordErr = errors.New("неправильный пароль")
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Password string `json:"password"`
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		utils.WriteBadRequestError(w)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &input); err != nil {
		utils.WriteBadRequestError(w)
		return
	}

	envPassword := utils.GetEnvOrDefault(passwordKey, "")
	if envPassword != input.Password {
		utils.WriteFailedValidationError(w, IncorrectPasswordErr)
		return
	}

	token, err := utils.GenerateJWTToken(envPassword)
	if err != nil {
		utils.WriteInternalServerError(w)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JsonData{"token": token}, nil)
}
