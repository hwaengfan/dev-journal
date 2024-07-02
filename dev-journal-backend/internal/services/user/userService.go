package userService

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/hwaengfan/dev-journal-backend/configs"
	userModel "github.com/hwaengfan/dev-journal-backend/internal/models/user"
	authenticationServices "github.com/hwaengfan/dev-journal-backend/internal/services/authentication"
	"github.com/hwaengfan/dev-journal-backend/internal/utils"
)

type Handler struct {
	store userModel.UserStore
}

func NewHandler(store userModel.UserStore) *Handler {
	return &Handler{store: store}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", handler.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", handler.handleRegister).Methods(http.MethodPost)
}

// Handler function for user login
func (handler *Handler) handleLogin(writer http.ResponseWriter, request *http.Request) {
	// get JSON payload
	var payload userModel.LoginUserPayload
	if error := utils.ParseJSON(request, &payload); error != nil {
		utils.WriteError(writer, http.StatusBadRequest, error)
		return
	}

	// validate payload
	if error := utils.Validate.Struct(payload); error != nil {
		errors := error.(validator.ValidationErrors)
		utils.WriteInvalidPayload(writer, errors)
		return
	}

	// validate user authentication
	user, error := handler.store.GetUserByEmail(payload.Email)
	if error != nil {
		utils.WriteError(writer, http.StatusNotFound, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !authenticationServices.ComparePassword(user.Password, []byte(payload.Password)) {
		utils.WriteError(writer, http.StatusNotFound, fmt.Errorf("not found, invalid email or password"))
		return
	}

	// create JWT token
	secret := []byte(configs.GlobalEnvironmentVariables.JWTSecret)
	token, error := authenticationServices.CreateJWT(secret, user.ID)
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, fmt.Errorf("failed to create JWT token: %v", error))
		return
	}

	utils.WriteJSON(writer, http.StatusOK, map[string]string{"token": token})
}

// Handler function for user registration
func (handler *Handler) handleRegister(writer http.ResponseWriter, request *http.Request) {
	// get JSON payload
	var payload userModel.RegisterUserPayload
	if error := utils.ParseJSON(request, &payload); error != nil {
		utils.WriteError(writer, http.StatusBadRequest, error)
		return
	}

	// validate payload
	if error := utils.Validate.Struct(payload); error != nil {
		errors := error.(validator.ValidationErrors)
		utils.WriteInvalidPayload(writer, errors)
		return
	}

	// check if user exists
	_, error := handler.store.GetUserByEmail(payload.Email)
	if error == nil {
		utils.WriteError(writer, http.StatusConflict, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	// hash password
	hashedPassword, error := authenticationServices.HashPassword(payload.Password)
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, error)
		return
	}

	// create user
	error = handler.store.CreateUser(userModel.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, error)
		return
	}

	utils.WriteJSON(writer, http.StatusCreated, nil)
}
