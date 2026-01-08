package tasks

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"go-final-project/pkg/api/nextdate"
	"go-final-project/pkg/db"
	"go-final-project/pkg/utils"
)

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {

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

	if task.Repeat != "" {
		nextDate, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			utils.WriteInternalServerError(w)
			return
		}

		err = db.UpdateTaskDate(nextDate, task.ID)
		if err != nil {
			utils.WriteInternalServerError(w)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.JsonData{}, nil)
		return
	}

	err = db.DeleteTask(taskId)
	if err != nil {
		utils.WriteInternalServerError(w)
		return
	} else {
		utils.WriteJSON(w, http.StatusOK, utils.JsonData{}, nil)
		return
	}
}
