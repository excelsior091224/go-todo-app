package handler

import (
	"fmt"
	"go-todo-app/model"
	"go-todo-app/repository"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// TaskIndex ...
func TaskIndex(c echo.Context) error {
	if c.Request().URL.Path == "/tasks" {
		c.Redirect(http.StatusPermanentRedirect, "/")
	}

	// タスクの一覧を取得する
	tasks, err := repository.TaskList()
	if err != nil {
		log.Println(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Tasks": tasks,
	}
	return render(c, "task/index.html", data)
}

// TaskNew ...
func TaskNew(c echo.Context) error {
	data := map[string]interface{}{
		"Message": "test",
		"Now":     time.Now(),
	}
	return render(c, "task/new.html", data)
}

// TaskShow ...
func TaskShow(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	task, err := repository.TaskGetByID(id)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}
	data := map[string]interface{}{
		"Task": task,
	}
	return render(c, "task/show.html", data)
}

// TaskCreateOutput ...
type TaskCreateOutput struct {
	Task             *model.Task
	Message          string
	ValidationErrors []string
}

// TaskCreate ...
func TaskCreate(c echo.Context) error {
	var task model.Task
	var out TaskCreateOutput
	if err := c.Bind(&task); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, out)
	}

	if err := c.Validate(&task); err != nil {
		c.Logger().Error(err.Error())

		out.ValidationErrors = task.ValidationErrors(err)

		return c.JSON(http.StatusUnprocessableEntity, out)
	}

	res, err := repository.TaskCreate(&task)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, out)
	}

	id, _ := res.LastInsertId()

	task.ID = int(id)

	out.Task = &task

	return c.JSON(http.StatusOK, out)
}

// TaskDelete ...
func TaskDelete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := repository.TaskDelete(id); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, "")
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("Task %d is deleted.", id))
}

// TaskEdit ...
func TaskEdit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	task, err := repository.TaskGetByID(id)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Task": task,
	}
	return render(c, "task/edit.html", data)
}

// TaskUpdateOutput ...
type TaskUpdateOutput struct {
	Task             *model.Task
	Message          string
	ValidationErrors []string
}

// TaskUpdate ...
func TaskUpdate(c echo.Context) error {
	ref := c.Request().Referer()
	refID := strings.Split(ref, "/")[4]
	reqID := c.Param("id")

	if reqID != refID {
		return c.JSON(http.StatusBadRequest, "")
	}

	var task model.Task
	var out TaskUpdateOutput
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, out)
	}

	if err := c.Validate(&task); err != nil {
		out.ValidationErrors = task.ValidationErrors(err)
		return c.JSON(http.StatusUnprocessableEntity, out)
	}

	taskID, _ := strconv.Atoi(reqID)
	task.ID = taskID
	_, err := repository.TaskUpdate(&task)

	if err != nil {
		out.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, out)
	}

	out.Task = &task
	return c.JSON(http.StatusOK, out)
}
