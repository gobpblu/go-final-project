package tasks

import (
	"net/http"

	"go-final-project/pkg/db"
	"go-final-project/pkg/utils"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	taskId := r.URL.Query().Get("id")
	if taskId == "" {
		utils.WriteBadRequestErrorWithMessage(w, "не указан идентификатор")
		return
	}

	err := db.DeleteTask(taskId)
	if err != nil {
		utils.WriteInternalServerError(w)
		return
	} else {
		utils.WriteJSON(w, http.StatusOK, utils.JsonData{}, nil)
		return
	}
}
