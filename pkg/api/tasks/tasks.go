package tasks

import (
	"go-final-project/pkg/db"
	"go-final-project/pkg/utils"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	searchQuery := query.Get("search")

	tasks, err := db.Tasks(50, searchQuery) // в параметре максимальное количество записей
	if err != nil {
		utils.WriteInternalServerError(w)
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}
	utils.WriteJSON(w, http.StatusOK, TasksResp{Tasks: tasks}, nil)
}
