package model

import (
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type Task struct {
	ID      int       `db:"id"`
	Title   string    `db:"title" validate:"required"`
	Text    string    `db:"text"`
	Status  int       `db:"status"`
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
}

// ValidationErrors ...
func (t *Task) ValidationErrors(err error) []string {
	var errMessages []string

	for _, err := range err.(validator.ValidationErrors) {
		var message string

		switch err.Field() {
		case "Title":
			message = "タイトルは必須です。"
			// case "Status":
			// 	message = "状態の選択は必須です。"
		}

		if message != "" {
			errMessages = append(errMessages, message)
		}
	}

	return errMessages
}
