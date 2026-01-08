package db

import (
	"database/sql"
	"errors"
	"time"

	"go-final-project/pkg/api/nextdate"
)

var (
	IncorrectIDErr = errors.New("неправильно указан идентификатор")
)

const (
	dateFormat = "02.01.2006"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(
		query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int, search string) ([]*Task, error) {
	var tasks []*Task
	var rows *sql.Rows

	parsedDate, err := time.Parse(dateFormat, search)

	if err == nil {
		date := parsedDate.Format(nextdate.DateLayout)
		rows, err = db.Query(
			"SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT :limit ",
			sql.Named("date", date),
			sql.Named("limit", limit),
		)
	} else if len(search) != 0 {
		searchParam := "%" + search + "%"
		rows, err = db.Query(
			"SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
			sql.Named("search", searchParam),
			sql.Named("limit", limit),
		)
	} else {
		rows, err = db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", limit))
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	task := &Task{}

	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	return task, err
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := db.Exec(
		query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return IncorrectIDErr
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return IncorrectIDErr
	}

	return nil
}

func UpdateTaskDate(nextDate string, id string) error {
	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	res, err := db.Exec(query, sql.Named("date", nextDate), sql.Named("id", id))

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return IncorrectIDErr
	}
	return nil
}
