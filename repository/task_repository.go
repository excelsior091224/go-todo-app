package repository

import (
	"database/sql"
	"go-todo-app/model"
	"time"
)

// TaskList ...
func TaskList() ([]*model.Task, error) {
	query := `SELECT * FROM tasks ORDER BY updated DESC;`

	var tasks []*model.Task
	if err := db.Select(&tasks, query); err != nil {
		return nil, err
	}

	return tasks, nil
}

// TaskGetByID ...
func TaskGetByID(id int) (*model.Task, error) {
	query := `SELECT *
	FROM tasks
	WHERE id = ?;`

	var task model.Task

	if err := db.Get(&task, query, id); err != nil {
		return nil, err
	}

	return &task, nil
}

// TaskCreate ...
func TaskCreate(task *model.Task) (sql.Result, error) {
	now := time.Now()

	task.Created = now
	task.Updated = now

	query := `INSERT INTO tasks (title, text, status, created, updated)
	VALUES (:title, :text, :status, :created, :updated);`

	tx := db.MustBegin()

	res, err := tx.NamedExec(query, task)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return res, nil
}

// TaskDelete ...
func TaskDelete(id int) error {
	query := "DELETE FROM tasks WHERE id = ?;"

	tx := db.MustBegin()

	if _, err := tx.Exec(query, id); err != nil {
		tx.Rollback()

		return err
	}
	return tx.Commit()
}

// TaskUpdate ...
func TaskUpdate(task *model.Task) (sql.Result, error) {
	now := time.Now()

	task.Updated = now

	query := `UPDATE tasks
	SET title = :title,
	text = :text,
	status = :status,
	updated = :updated
	WHERE id = :id;`

	tx := db.MustBegin()

	res, err := tx.NamedExec(query, task)

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return res, nil
}
