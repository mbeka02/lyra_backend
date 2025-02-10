package server

import (
	"fmt"
	"net/http"

	"github.com/mbeka02/lyra_backend/internal/model"
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
	request := model.CreateUserRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	response, err := h.userService.CreateUser(r.Context(), request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	if err := respondWithJSON(w, http.StatusCreated, response); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	request := model.UpdateUserRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	// ensure auth payload is present
	payload, err := getAuthPayload(r.Context())
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

	if err := respondWithJSON(w, http.StatusCreated, "account updated"); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *UserHandler) HandleProfilePicture(w http.ResponseWriter, r *http.Request) {
	// get the file
	_, fileHeader, err := r.FormFile("image")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("bad request:%v", err))
		return
	}
	// ensure auth payload is present
	payload, err := getAuthPayload(r.Context())
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

	if err := respondWithJSON(w, http.StatusCreated, "profile picture updated"); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
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

	if err := respondWithJSON(w, http.StatusOK, response); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}
