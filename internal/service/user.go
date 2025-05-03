package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo"
	"family-flow-app/internal/repo/repoerrs"
	"family-flow-app/pkg/hasher"
	"family-flow-app/pkg/token"
)

type UserService struct {
	userRepo repo.User
}

func NewUserService(userRepo repo.User) *UserService {
	return &UserService{userRepo: userRepo}
}

func (u *UserService) Login(ctx context.Context, log *slog.Logger, input AuthInput) (string, error) {
	log.Info(fmt.Sprintf("Service - UserService - Login"))
	var err error
	var tokenString string
	var output entity.User
	output, err = u.isExist(ctx, log, input)

	if errors.Is(err, ErrInvalidPassword) {
		return "", ErrInvalidPassword
	} else if errors.Is(err, ErrUserNotFound) {
		return "", ErrUserNotFound
	} else if err == nil {
		if tokenString, err = token.Create(output.Id); err != nil {
			return "", err
		}
		return tokenString, err
	} else {
		return "", err
	}
}

func (u *UserService) Create(ctx context.Context, log *slog.Logger, input UserCreateInput) (string, error) {
	var err error
	var tokenString string
	log.Info(fmt.Sprintf("Service - UserService - Create"))
	//hash
	password, err := hasher.HashPassword(input.Password)
	if err != nil {
		return "", ErrCannotHashPassword
	}
	user := entity.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: password,
		Role:     input.Role,
	}

	id, err := u.userRepo.Create(ctx, user)
	if err != nil {
		if err == repoerrs.ErrAlreadyExists {
			return "", ErrUserAlreadyExists
		}
		log.Error(fmt.Sprintf("Service - UserService - Create: %v", err))
		return "", ErrCannotCreateUser
	}
	log.Info(fmt.Sprintf("Service - UserService - userRepo.Create - id: %s", id))

	// token
	if tokenString, err = token.Create(id); err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *UserService) GetById(ctx context.Context, log *slog.Logger, id string) (entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.User{}, ErrUserNotFound
		}
		log.Error(fmt.Sprintf("Service - UserService - GetById: %v", err))
		return entity.User{}, ErrCannotGetUser
	}
	return user, nil
}

func (u *UserService) GetByEmail(ctx context.Context, log *slog.Logger, input UserGetByEmailInput) (
	entity.User, error,
) {
	user, err := u.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.User{}, ErrUserNotFound
		}
		log.Error(fmt.Sprintf("Service - UserService - GetByEmail: %v", err))
		return entity.User{}, ErrCannotGetUser
	}
	return user, nil
}

func (u *UserService) isExist(ctx context.Context, log *slog.Logger, input AuthInput) (entity.User, error) {
	var err error
	log.Info(fmt.Sprintf("Service - UserService - isExist"))
	output, err := u.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		log.Error(fmt.Sprintf("Service - UserService - isExist - GetByEmail: %v", err))
		return entity.User{}, ErrUserNotFound
	}

	if err = hasher.CheckPassword(input.Password, output.Password); err != nil {
		log.Error(fmt.Sprintf("Service - UserService - isExist - CheckPassword: %v", err))
		return entity.User{}, ErrInvalidPassword
	}

	return output, err
}

func (u *UserService) AddMemberToFamily(ctx context.Context, log *slog.Logger, input AddMemberToFamilyInput) error {
	log.Info("Service - UserService - AddMember")
	var err error

	// update role
	err = u.userRepo.UpdateRole(ctx, input.UserEmail, input.Role)
	if err != nil {
		log.Error(fmt.Sprintf("Service - UserService - AddMember: %v", err))
		return ErrCannotUpdateUser
	}

	err = u.userRepo.UpdateFamilyID(ctx, input.FamilyId, input.FamilyId)
	if err != nil {
		log.Error(fmt.Sprintf("Service - UserService - AddMember: %v", err))
		return ErrCannotAddMemberToFamily
	}

	return nil
}

type UpdateUserInput struct {
	ID    string
	Name  string
	Email string
	Role  string
}

func (u *UserService) Update(ctx context.Context, log *slog.Logger, input UpdateUserInput) error {
	log.Info(fmt.Sprintf("Service - UserService - Update"))
	user := entity.User{
		Id:    input.ID,
		Name:  input.Name,
		Email: input.Email,
		Role:  input.Role,
	}

	err := u.userRepo.Update(ctx, user)
	if err != nil {
		log.Error(fmt.Sprintf("Service - UserService - Update: %v", err))
		return ErrCannotUpdateUser
	}
	return nil
}

func (u *UserService) UpdatePassword(ctx context.Context, log *slog.Logger, email, password string) error {
	log.Info(fmt.Sprintf("Service - UserService - UpdatePassword"))
	//hash
	passwordHash, err := hasher.HashPassword(password)
	if err != nil {
		return ErrCannotHashPassword
	}

	err = u.userRepo.UpdatePassword(ctx, email, passwordHash)
	if err != nil {
		log.Error(fmt.Sprintf("Service - UserService - UpdatePassword: %v", err))
		return ErrCannotUpdateUser
	}
	return nil
}

func (u *UserService) ResetFamilyID(ctx context.Context, log *slog.Logger, id string) error {
	log.Info(fmt.Sprintf("Service - UserService - ResetFamilyID"))

	// Получение текущего пользователя из контекста
	currentUser, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("Service - UserService - ResetFamilyID - GetCurrentUser: %v", err))
		return ErrCannotResetFamilyID
	}

	// Проверка роли пользователя
	if currentUser.Role != "Parent" && currentUser.Id != id {
		log.Error("Service - UserService - ResetFamilyID: insufficient permissions")
		return ErrInsufficientPermissions
	}

	// Сброс FamilyID
	err = u.userRepo.ResetFamilyID(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("Service - UserService - ResetFamilyID: %v", err))
		return ErrCannotResetFamilyID
	}
	return nil
}

func (u *UserService) ExistsByEmail(ctx context.Context, log *slog.Logger, email string) (bool, error) {
	log.Info("Service - UserService - ExistsByEmail")
	return u.userRepo.ExistsByEmail(ctx, email)
}
