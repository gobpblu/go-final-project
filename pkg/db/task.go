package db

import (
	"database/sql"
	"go-final-project/pkg/api/nextdate"
	"time"
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
	// определите запрос
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
			"SELECT * FROM scheduler WHERE date = :date LIMIT :limit ",
			sql.Named("date", date),
			sql.Named("limit", limit),
		)
	} else if len(search) != 0 {
		searchParam := "%" + search + "%"
		rows, err = db.Query(
			"SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
			sql.Named("search", searchParam),
			sql.Named("limit", limit),
		)
	} else {
		rows, err = db.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", limit))
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
