package noteService

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	noteModel "github.com/hwaengfan/dev-journal-backend/internal/models/note"
	projectModel "github.com/hwaengfan/dev-journal-backend/internal/models/project"
	userModel "github.com/hwaengfan/dev-journal-backend/internal/models/user"
	authenticationServices "github.com/hwaengfan/dev-journal-backend/internal/services/authentication"
	"github.com/hwaengfan/dev-journal-backend/internal/utils"
)

type Handler struct {
	store        noteModel.NoteStore
	userStore    userModel.UserStore
	projectStore projectModel.ProjectStore
}

func NewHandler(store noteModel.NoteStore, userStore userModel.UserStore, projectStore projectModel.ProjectStore) *Handler {
	return &Handler{store: store, userStore: userStore, projectStore: projectStore}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/notes/create-new-note", authenticationServices.JWTAuthentication(handler.handleCreateNewNote, handler.userStore)).Methods(http.MethodPost)

	router.HandleFunc("/notes/get-notes-by-linked-project-ID/{projectID}", authenticationServices.JWTAuthentication(handler.handleGetNotesByLinkedProjectID, handler.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/notes/get-notes-by-ID/{noteID}", authenticationServices.JWTAuthentication(handler.handleGetNoteByID, handler.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/notes/update-note-by-ID/{noteID}", authenticationServices.JWTAuthentication(handler.handleUpdateNoteByID, handler.userStore)).Methods(http.MethodPut)

	router.HandleFunc("/notes/delete-note-by-ID/{noteID}", authenticationServices.JWTAuthentication(handler.handleDeleteNoteByID, handler.userStore)).Methods(http.MethodDelete)
}

// Handler function for creating a new note
func (handler *Handler) handleCreateNewNote(writer http.ResponseWriter, request *http.Request) {
	// validate if the user is logged in
	userID := authenticationServices.GetUserIDFromContext(request.Context())
	if !userID.Valid {
		utils.WritePermissionDenied(writer)
		return
	}

	// get JSON payload
	var payload noteModel.CreateNotePayload
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

	// check if the linkedProjectID is provided
	if payload.LinkedProjectID == uuid.Nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("missing linked project ID"))
		return
	}

	// check if the project exists
	if error := handler.validateLinkedProjectID(payload.LinkedProjectID); error != nil {
		utils.WriteError(writer, http.StatusBadRequest, error)
		return
	}

	// insert the new note into the database
	noteID, error := handler.store.CreateNote(noteModel.Note{
		UserID:          userID.UUID,
		LinkedProjectID: payload.LinkedProjectID,
		Title:           payload.Title,
		Content:         payload.Content,
		Favorited:       payload.Favorited,
		Tags:            payload.Tags,
	})
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, error)
		return
	}

	utils.WriteJSON(writer, http.StatusCreated, map[string]string{"noteID": noteID.String()})
}

// Handler function for getting notes by linked project ID
func (handler *Handler) handleGetNotesByLinkedProjectID(writer http.ResponseWriter, request *http.Request) {
	// validate if the user is logged in
	userID := authenticationServices.GetUserIDFromContext(request.Context())
	if !userID.Valid {
		utils.WritePermissionDenied(writer)
		return
	}

	// get project ID from URL
	linkedProjectIDString, exists := mux.Vars(request)["projectID"]
	if !exists {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("missing linked project ID in parameters"))
		return
	}

	linkedProjectID, error := uuid.Parse(linkedProjectIDString)
	if error != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid linked project ID"))
		return
	}

	// check if the linkedProjectID is provided
	if linkedProjectID == uuid.Nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("null linked project ID"))
		return
	}

	// check if the project exists
	if error := handler.validateLinkedProjectID(linkedProjectID); error != nil {
		utils.WriteError(writer, http.StatusBadRequest, error)
		return
	}

	// get notes by projectID
	notes, error := handler.store.GetNotesByLinkedProjectID(linkedProjectID)
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, error)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, notes)
}

// Handler function for getting a note by ID
func (handler *Handler) handleGetNoteByID(writer http.ResponseWriter, request *http.Request) {
	// validate if the user is logged in
	userID := authenticationServices.GetUserIDFromContext(request.Context())
	if !userID.Valid {
		utils.WritePermissionDenied(writer)
		return
	}

	// get noteID from URL
	noteIDString, exists := mux.Vars(request)["noteID"]
	if !exists {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("missing note ID"))
		return
	}

	noteID, error := uuid.Parse(noteIDString)
	if error != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid note ID"))
		return
	}

	note, error := handler.store.GetNoteByID(noteID)
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, fmt.Errorf("failed to get note by ID: %v", error))
		return
	}

	utils.WriteJSON(writer, http.StatusOK, note)
}

// Handler function for updating a note by ID
func (handler *Handler) handleUpdateNoteByID(writer http.ResponseWriter, request *http.Request) {
	// validate if the user is logged in
	userID := authenticationServices.GetUserIDFromContext(request.Context())
	if !userID.Valid {
		utils.WritePermissionDenied(writer)
		return
	}

	// get noteID from URL
	noteIDString, exists := mux.Vars(request)["noteID"]
	if !exists {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("missing note ID"))
		return
	}

	noteID, error := uuid.Parse(noteIDString)
	if error != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid note ID"))
		return
	}

	// get JSON payload
	var payload noteModel.UpdateNotePayload
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

	// check if the note exists
	error = handler.validateNoteID(noteID)
	if error != nil {
		utils.WriteError(writer, http.StatusBadRequest, error)
		return
	}

	// check if the linkedProjectID is provided
	if payload.LinkedProjectID != uuid.Nil {
		// check if the project exists
		if error := handler.validateLinkedProjectID(payload.LinkedProjectID); error != nil {
			utils.WriteError(writer, http.StatusBadRequest, error)
			return
		}
	}

	// update the note by ID
	error = handler.store.UpdateNoteByID(noteModel.Note{
		LinkedProjectID: payload.LinkedProjectID,
		Title:           payload.Title,
		Content:         payload.Content,
		Favorited:       payload.Favorited,
		Tags:            payload.Tags,
	}, noteID)
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

// Handler function for deleting a note by ID
func (handler *Handler) handleDeleteNoteByID(writer http.ResponseWriter, request *http.Request) {
	// validate if the user is logged in
	userID := authenticationServices.GetUserIDFromContext(request.Context())
	if !userID.Valid {
		utils.WritePermissionDenied(writer)
		return
	}

	// get noteID from URL
	noteIDString, exists := mux.Vars(request)["noteID"]
	if !exists {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("missing note ID"))
		return
	}

	noteID, error := uuid.Parse(noteIDString)
	if error != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid note ID"))
		return
	}

	// check if the note exists
	_, error = handler.store.GetNoteByID(noteID)
	if error != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("note does not exist"))
		return
	}

	// delete the note by ID
	error = handler.store.DeleteNoteByID(noteID)
	if error != nil {
		utils.WriteError(writer, http.StatusInternalServerError, error)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, nil)
}

// validateLinkedProjectID check if the project exists
func (handler *Handler) validateLinkedProjectID(linkedProjectID uuid.UUID) error {
	_, error := handler.projectStore.GetProjectByID(linkedProjectID)
	if error != nil {
		return fmt.Errorf("project ID does not exist")
	}

	return nil
}

// validateNoteID check if the note exists
func (handler *Handler) validateNoteID(noteID uuid.UUID) error {
	_, error := handler.store.GetNoteByID(noteID)
	if error != nil {
		return fmt.Errorf("note ID does not exist")
	}

	return nil
}
