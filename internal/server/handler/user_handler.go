package handler

import (
	"fmt"
	"net/http"

	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/middleware"
	"github.com/mbeka02/lyra_backend/internal/server/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var request model.CreateUserRequest
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)

		return
	}

	user, err := h.userService.CreateUser(r.Context(), request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	user, err := h.userService.GetUser(r.Context(), payload.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var request model.UpdateUserRequest
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	// ensure auth payload is present
	payload, err := middleware.GetAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	// update the user account
	err = h.userService.UpdateUser(r.Context(), model.UpdateUserRequest{
		Email:           request.Email,
		TelephoneNumber: request.TelephoneNumber,
		FullName:        request.FullName,
	}, payload.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, "account updated")
}

func (h *UserHandler) HandleProfilePicture(w http.ResponseWriter, r *http.Request) {
	// get the file
	_, fileHeader, err := r.FormFile("image")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("bad request:%v", err))
		return
	}
	// ensure auth payload is present
	payload, err := middleware.GetAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	// update the profile picture
	err = h.userService.UpdateProfilePicture(r.Context(), fileHeader, payload.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, "profile picture updated")
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	request := model.LoginRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	response, err := h.userService.Login(r.Context(), request)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err)
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}
