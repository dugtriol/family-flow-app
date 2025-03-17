package v1

import "fmt"

var (
	MsgInvalidReq        = "Invalid request"
	MsgFailedParsing     = "Failed to parse data"
	MsgInternalServerErr = "Internal server error"

	MsgUserNotFound      = "User not found"
	MsgUserAlreadyExists = "User already exists"

	MsgInvalidPasswordErr = "Invalid password"
	ErrInvalidToken       = fmt.Errorf("invalid token")
	ErrUserGet            = fmt.Errorf("user not get from database")
	ErrNoUserInContext    = fmt.Errorf("no user in the context")

	MsgTaskNotFound = "Task not found"
)
