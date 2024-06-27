package userService

import (
	"fmt"
	"net/http"
	
	"github.com/gorilla/mux"
	"github.com/hwaengfan/dev-journal-backend/internal/utils"
	"github.com/hwaengfan/dev-journal-backend/internal/models/user"
	"github.com/hwaengfan/dev-journal-backend/internal/services/authentication"
)

type Handler struct {
	store userModel.UserStore
}

func NewHandler(store userModel.UserStore) *Handler {
	return &Handler{store: store}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", handler.handleLogin).Methods("POST")
	router.HandleFunc("/register", handler.handleRegister).Methods("POST")
}

func (handler *Handler) handleLogin(writer http.ResponseWriter, request *http.Request) {
	// handle login
}

func (handler *Handler) handleRegister(writer http.ResponseWriter, request *http.Request) {
	// get JSON payload
	var payload userModel.RegisterUserPayload
	if err := utils.ParseJSON(request, payload); err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
	}

	// check if user exists
	_, error := handler.store.GetUserByEmail(payload.Email)
	if error == nil {
		utils.WriteError(writer, http.StatusConflict, fmt.Errorf("User with email %s already exists", payload.Email))
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
		LastName: payload.LastName,
		Email: payload.Email,
		Password: hashedPassword,
	})

	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, error)
		return
	}

	utils.WriteJSON(writer, http.StatusCreated, nil)
}
