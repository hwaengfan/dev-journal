package taskService

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
	store        taskModel.TaskStore
	userStore    userModel.UserStore
	projectStore projectModel.ProjectStore
}

func NewHandler(store taskModel.TaskStore, userStore userModel.UserStore, projectStore projectModel.ProjectStore) *Handler {
	return &Handler{store: store, userStore: userStore, projectStore: projectStore}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/tasks/create-new-task", authenticationServices.JWTAuthentication(handler.handleCreateNewTask, handler.userStore)).Methods(http.MethodPost)

	router.HandleFunc("/tasks/get-tasks-by-linked-project-ID/{projectID}", authenticationServices.JWTAuthentication(handler.handleGetTasksByLinkedProjectID, handler.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/tasks/update-task-by-ID/{taskID}", authenticationServices.JWTAuthentication(handler.handleUpdateTaskByID, handler.userStore)).Methods(http.MethodPut)

	router.HandleFunc("/tasks/delete-task-by-ID/{taskID}", authenticationServices.JWTAuthentication(handler.handleDeleteTaskByID, handler.userStore)).Methods(http.MethodDelete)
}

// Handler function for creating a new task
func (handler *Handler) handleCreateNewTask(writer http.ResponseWriter, request *http.Request) {
	// validate if the user is logged in
	userID := authenticationServices.GetUserIDFromContext(request.Context())
	if !userID.Valid {
		utils.WritePermissionDenied(writer)
		return
	}

	// get JSON payload
	var payload taskModel.CreateTaskPayload
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

	// check if the project exists
	_, error := handler.projectStore.GetProjectByID(payload.LinkedProjectID)
	if error != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("project does not exist to link task to"))
		return
	}

	// insert the new task into the database
	taskID, error := handler.store.CreateTask(taskModel.Task{
		LinkedProjectID: payload.LinkedProjectID,
		Description:     payload.Description,
		Completed:       payload.Completed,
	})
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, error)
		return
	}

	utils.WriteJSON(writer, http.StatusCreated, map[string]uuid.UUID{"taskID": taskID})
}

// Handler function for getting all tasks in a project
func (handler *Handler) handleGetTasksByLinkedProjectID(writer http.ResponseWriter, request *http.Request) {
	// validate if the user is logged in
	userID := authenticationServices.GetUserIDFromContext(request.Context())
	if !userID.Valid {
		utils.WritePermissionDenied(writer)
		return
	}

	// get projectID from URL
	linkedProjectIDString, exists := mux.Vars(request)["projectID"]
	if !exists {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("missing project ID"))
		return
	}

	linkedProjectID, error := uuid.Parse(linkedProjectIDString)
	if error != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid project ID"))
		return
	}

	// get tasks by projectID
	tasks, error := handler.store.GetTasksByLinkedProjectID(linkedProjectID)
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, error)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, tasks)
}

// Handler function for updating a task by ID
func (handler *Handler) handleUpdateTaskByID(writer http.ResponseWriter, request *http.Request) {
	// validate if the user is logged in
	userID := authenticationServices.GetUserIDFromContext(request.Context())
	if !userID.Valid {
		utils.WritePermissionDenied(writer)
		return
	}

	// get taskID from URL
	taskIDString, exists := mux.Vars(request)["taskID"]
	if !exists {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("missing task ID"))
		return
	}

	taskID, error := uuid.Parse(taskIDString)
	if error != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid task ID"))
		return
	}

	// get JSON payload
	var payload taskModel.UpdateTaskPayload
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

	// check if the project exists if the linkedProjectID is provided
	if payload.LinkedProjectID != uuid.Nil {
		_, error := handler.projectStore.GetProjectByID(payload.LinkedProjectID)
		if error != nil {
			utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("project does not exist to link task to"))
			return
		}
	}

	// update the task by ID
	error = handler.store.UpdateTaskByID(taskModel.Task{
		LinkedProjectID: payload.LinkedProjectID,
		Description:     payload.Description,
		Completed:       payload.Completed,
	}, taskID)
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

// Handler function for deleting a task by ID
func (handler *Handler) handleDeleteTaskByID(writer http.ResponseWriter, request *http.Request) {
	// validate if the user is logged in
	userID := authenticationServices.GetUserIDFromContext(request.Context())
	if !userID.Valid {
		utils.WritePermissionDenied(writer)
		return
	}

	// get taskID from URL
	taskIDString, exists := mux.Vars(request)["taskID"]
	if !exists {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("missing task ID"))
		return
	}

	taskID, error := uuid.Parse(taskIDString)
	if error != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid task ID"))
		return
	}

	// check if the task exists
	_, error = handler.store.GetTaskByID(taskID)
	if error != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("task does not exist"))
		return
	}

	// delete the task by ID
	error = handler.store.DeleteTaskByID(taskID)
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, error)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, nil)
}
