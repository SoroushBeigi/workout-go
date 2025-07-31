package api

import (
	"errors"
	"log"
	"regexp"

	"github.com/SoroushBeigi/workout-go/internal/store"
)

type RegisteUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) validateRegisterRequest(req *RegisteUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) > 50 {
		return errors.New("username is too long")
	}
	if len(req.Username) < 4 {
		return errors.New("username is too short")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) > 50 {
		return errors.New("password is too long")
	}
	if len(req.Password) < 8 {
		return errors.New("password is too short")
	}

	return nil
}
