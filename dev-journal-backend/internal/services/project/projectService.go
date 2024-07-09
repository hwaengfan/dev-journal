package projectService

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	projectModel "github.com/hwaengfan/dev-journal-backend/internal/models/project"
	taskModel "github.com/hwaengfan/dev-journal-backend/internal/models/task"
	userModel "github.com/hwaengfan/dev-journal-backend/internal/models/user"
	authenticationServices "github.com/hwaengfan/dev-journal-backend/internal/services/authentication"
	"github.com/hwaengfan/dev-journal-backend/internal/utils"
)

type Handler struct {
	store     projectModel.ProjectStore
	userStore userModel.UserStore
	taskStore taskModel.TaskStore
}

func NewHandler(store projectModel.ProjectStore, userStore userModel.UserStore, taskStore taskModel.TaskStore) *Handler {
	return &Handler{store: store, userStore: userStore, taskStore: taskStore}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/projects/create-new-project", authenticationServices.JWTAuthentication(handler.handleCreateNewProject, handler.userStore)).Methods(http.MethodPost)

	router.HandleFunc("/projects/get-projects-by-user-ID", authenticationServices.JWTAuthentication(handler.handleGetProjectsByUserID, handler.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/projects/get-project-by-ID/{projectID}", authenticationServices.JWTAuthentication(handler.handleGetProjectByID, handler.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/projects/update-project-by-ID/{projectID}", authenticationServices.JWTAuthentication(handler.handleUpdateProjectByID, handler.userStore)).Methods(http.MethodPut)

	router.HandleFunc("/projects/delete-project-by-ID/{projectID}", authenticationServices.JWTAuthentication(handler.handleDeleteProjectByID, handler.userStore)).Methods(http.MethodDelete)
}

// Handler function for creating a new project
func (handler *Handler) handleCreateNewProject(writer http.ResponseWriter, request *http.Request) {
	// validate if the user is logged in
	userID := authenticationServices.GetUserIDFromContext(request.Context())
	if !userID.Valid {
		utils.WritePermissionDenied(writer)
		return
	}

	// get JSON payload
	var payload projectModel.CreateProjectPayload
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

	// insert the new project into the database
	projectID, error := handler.store.CreateProject(projectModel.Project{
		UserID:      userID.UUID,
		Title:       payload.Title,
		Description: payload.Description,
		Priority:    payload.Priority,
		Deadline:    payload.Deadline,
	})
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, error)
		return
	}

	utils.WriteJSON(writer, http.StatusCreated, map[string]uuid.UUID{"projectID": projectID})
}

// Handler function for getting all projects by user ID
func (handler *Handler) handleGetProjectsByUserID(writer http.ResponseWriter, request *http.Request) {
	// validate if the user is logged in
	userID := authenticationServices.GetUserIDFromContext(request.Context())
	if !userID.Valid {
		utils.WritePermissionDenied(writer)
		return
	}

	// get the user's projects
	projects, error := handler.store.GetProjectsByUserID(userID.UUID)
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, error)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, projects)
}

// Handler function for getting a project by ID
func (handler *Handler) handleGetProjectByID(writer http.ResponseWriter, request *http.Request) {
	// validate if the user is logged in
	userID := authenticationServices.GetUserIDFromContext(request.Context())
	if !userID.Valid {
		utils.WritePermissionDenied(writer)
		return
	}

	// get project ID from URL
	projectIDString, exists := mux.Vars(request)["projectID"]
	if !exists {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("missing project ID"))
		return
	}

	projectID, error := uuid.Parse(projectIDString)
	if error != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid project ID"))
		return
	}

	// get the project by ID
	project, error := handler.store.GetProjectByID(projectID)
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, fmt.Errorf("failed to get project by ID: %v", error))
		return
	}

	utils.WriteJSON(writer, http.StatusOK, project)
}

// Handler function for updating a project by ID
func (handler *Handler) handleUpdateProjectByID(writer http.ResponseWriter, request *http.Request) {
	// validate if the user is logged in
	userID := authenticationServices.GetUserIDFromContext(request.Context())
	if !userID.Valid {
		utils.WritePermissionDenied(writer)
		return
	}

	// get project ID from URL
	projectIDString, exists := mux.Vars(request)["projectID"]
	if !exists {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("missing project ID"))
		return
	}

	projectID, error := uuid.Parse(projectIDString)
	if error != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid project ID"))
		return
	}

	// get JSON payload
	var payload projectModel.UpdateProjectPayload
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

	// update the project by ID
	error = handler.store.UpdateProjectByID(projectModel.Project{
		Title:       payload.Title,
		Description: payload.Description,
		Priority:    payload.Priority,
		Deadline:    payload.Deadline,
	}, projectID)
	if error != nil {
		if error.Error() == "no fields to update" {
			utils.WriteError(writer, http.StatusBadRequest, error)
		} else {
			utils.WriteError(writer, http.StatusInternalServerError, error)
		}
		return
	}

	utils.WriteJSON(writer, http.StatusOK, nil)
}

// Handler function for deleting a project by ID
func (handler *Handler) handleDeleteProjectByID(writer http.ResponseWriter, request *http.Request) {
	// validate if the user is logged in
	userID := authenticationServices.GetUserIDFromContext(request.Context())
	if !userID.Valid {
		utils.WritePermissionDenied(writer)
		return
	}

	// get project ID from URL
	projectIDString, exists := mux.Vars(request)["projectID"]
	if !exists {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("missing project ID"))
		return
	}

	projectID, error := uuid.Parse(projectIDString)
	if error != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid project ID"))
		return
	}

	// delete all tasks linked to the project by ID
	error = handler.taskStore.DeleteTasksByLinkedProjectID(projectID)
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, error)
		return
	}

	// delete the project by ID
	error = handler.store.DeleteProjectByID(projectID)
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, error)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, nil)
}
