package service

import "fmt"

var (
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrCannotCreateUser  = fmt.Errorf("cannot create user")
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrCannotGetUser     = fmt.Errorf("cannot get user")

	ErrCannotHashPassword = fmt.Errorf("cannot hash password")
	ErrInvalidPassword    = fmt.Errorf("invalid password")

	ErrForbidden = fmt.Errorf("forbidden")
	ErrCode      = fmt.Errorf("err code")

	ErrCannotCreateFamily      = fmt.Errorf("cannot create family")
	ErrCannotAddMemberToFamily = fmt.Errorf("cannot add member to family")
	ErrFamilyNotFound          = fmt.Errorf("family not found")
	ErrCannotGetFamilyMembers  = fmt.Errorf("cannot get family members")
	ErrCannotCreateList        = fmt.Errorf("cannot create list")
	ErrListNotFound            = fmt.Errorf("list not found")
	ErrCannotUpdateList        = fmt.Errorf("cannot update list")
	ErrCannotDeleteList        = fmt.Errorf("cannot delete list")

	ErrCannotCreateItem = fmt.Errorf("cannot create item")
	ErrItemNotFound     = fmt.Errorf("item not found")
	ErrCannotUpdateItem = fmt.Errorf("cannot update item")
	ErrCannotDeleteItem = fmt.Errorf("cannot delete item")
)
