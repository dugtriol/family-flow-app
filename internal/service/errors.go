package service

import "fmt"

var (
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrCannotCreateUser  = fmt.Errorf("cannot create user")
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrCannotGetUser     = fmt.Errorf("cannot get user")

	ErrCannotHashPassword = fmt.Errorf("cannot hash password")
	ErrInvalidPassword    = fmt.Errorf("invalid password")

	ErrTaskAlreadyExists = fmt.Errorf("task already exists")
	ErrTaskNotFound      = fmt.Errorf("task not found")
	ErrCannotCreateTask  = fmt.Errorf("cannot create task")
	ErrCannotGetTask     = fmt.Errorf("cannot get task")
	ErrCannotGetTasks    = fmt.Errorf("cannot get tasks")
	ErrCannotUpdateTask  = fmt.Errorf("cannot update task")
	ErrCannotDeleteTask  = fmt.Errorf("cannot delete task")

	ErrCannotCompleteTask = fmt.Errorf("cannot complete task")
	ErrForbidden          = fmt.Errorf("forbidden")
)
