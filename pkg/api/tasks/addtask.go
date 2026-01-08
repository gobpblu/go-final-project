package tasks

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go-final-project/pkg/api/nextdate"
	"go-final-project/pkg/db"
	"go-final-project/pkg/utils"
)

type AddTaskResp struct {
	id int64
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string `json:"title"`
		Date    string `json:"date"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
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

	task := &db.Task{
		Title:   input.Title,
		Date:    input.Date,
		Comment: input.Comment,
		Repeat:  input.Repeat,
	}

	err = validateTask(task)
	if err != nil {
		utils.WriteFailedValidationError(w, err)
		return
	}

	id, err := db.AddTask(task)
	if err != nil {
		utils.WriteInternalServerError(w)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JsonData{"id": id}, nil)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		Date    string `json:"date"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
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

	task := &db.Task{
		ID:      input.ID,
		Title:   input.Title,
		Date:    input.Date,
		Comment: input.Comment,
		Repeat:  input.Repeat,
	}

	err = validateTask(task)
	if err != nil {
		utils.WriteFailedValidationError(w, err)
		return
	}

	err = db.UpdateTask(task)
	if err != nil {
		if errors.Is(err, db.IncorrectIDErr) {
			utils.WriteBadRequestErrorWithMessage(w, err.Error())
			return
		}
		utils.WriteInternalServerError(w)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JsonData{}, nil)
}

func validateTask(task *db.Task) error {
	if task.Title == "" {
		return errors.New("Поле title является обязательным")
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(nextdate.DateLayout)
	} else {
		date, err := time.Parse(nextdate.DateLayout, task.Date)
		if err != nil {
			return nextdate.InvalidDateParamErr
		}

		if nextdate.AfterNow(now, date) {
			if task.Repeat == "" {
				task.Date = now.Format(nextdate.DateLayout)
			} else {
				nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return err
				}
				task.Date = nextDate
			}
		}
	}

	return nil
}
