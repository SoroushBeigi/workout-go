package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/SoroushBeigi/workout-go/internal/store"
	"github.com/SoroushBeigi/workout-go/internal/utils"
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

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisteUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		utils.WriteJson(
			w, http.StatusBadRequest,
			utils.Envelope{"error": "Invalid request body"},
		)
		return
	}
	err = h.validateRegisterRequest(&req)
	if err != nil {
		utils.WriteJson(
			w, http.StatusBadRequest,
			utils.Envelope{"error": err.Error()},
		)
		return
	}
	user := store.User{
		Username: req.Username,
		Email:    req.Email,
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("Error setting password hash: %v", err)
		utils.WriteJson(
			w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"},
		)
		return
	}
	err = h.userStore.CreateUser(&user)
	if err != nil {
		h.logger.Printf("Error: registering user: %v", err)
		utils.WriteJson(
			w, http.StatusInternalServerError,
			utils.Envelope{
				"error": "Internal server error",
			},
		)
		return
	}
	utils.WriteJson(
		w, http.StatusCreated,
		utils.Envelope{
			"user":    user,
			"message": "User registered successfully",
		},
	)

}
