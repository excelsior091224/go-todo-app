package main

import (
	"log"

	"go-todo-app/handler"
	"go-todo-app/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/go-playground/validator.v9"
)

const dbPath = "db/db.sql"

var db *sqlx.DB
var e = createMux()

func main() {
	db = connectDB()
	repository.SetDB(db)

	e.GET("/", handler.TaskIndex)
	e.GET("/tasks", handler.TaskIndex)
	e.GET("/tasks/new", handler.TaskNew)
	e.GET("/tasks/:id", handler.TaskShow)
	e.GET("/tasks/:id/edit", handler.TaskEdit)
	e.POST("/api/tasks", handler.TaskCreate)
	e.DELETE("/api/tasks/:id", handler.TaskDelete)
	e.PATCH("/api/tasks/:id", handler.TaskUpdate)

	e.Logger.Fatal(e.Start(":8080"))
}

func createMux() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.Use(middleware.CSRF())

	e.Static("/css", "src/css")
	e.Static("/js", "src/js")

	e.Validator = &CustomValidator{validator: validator.New()}

	return e
}

func connectDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		e.Logger.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("db connection succeeded")
	return db
}

// CustomValidator ...
type CustomValidator struct {
	validator *validator.Validate
}

// Validate ...
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
