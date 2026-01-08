package tasks

import (
	"database/sql"
	"errors"
	"net/http"

	"go-final-project/pkg/db"
	"go-final-project/pkg/utils"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

	taskId := r.URL.Query().Get("id")
	if taskId == "" {
		utils.WriteBadRequestErrorWithMessage(w, "не указан идентификатор")
		return
	}

	task, err := db.GetTask(taskId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WritePageNotFoundError(w, "задача с указанным идентификатором не найдена")
			return
		}
		utils.WriteInternalServerError(w)
		return
	}

	utils.WriteJSON(w, http.StatusOK, task, nil)
}
