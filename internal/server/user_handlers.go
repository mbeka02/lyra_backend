package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/mbeka02/lyra_backend/internal/auth"
	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/models"
)

type UserHandler struct {
	Store               *database.Store
	AuthMaker           auth.Maker
	AccessTokenDuration time.Duration
}

func (uh *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	request := models.CreateUserRequest{}

	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	passwordHash, err := auth.HashPassword(request.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("failed to process password"))
		return
	}

	user, err := uh.Store.CreateUser(r.Context(), database.CreateUserParams{
		FullName: request.Fullname,
		Email:    request.Email,
		Password: passwordHash,
	})
	if err != nil {

		respondWithError(w, http.StatusInternalServerError, errors.New("failed to create user"))
		return
	}
	// response := APIResponse{
	// 	Status:  http.StatusCreated,
	// 	Message: "User created successfully",
	// 	Data:    user,
	// }
	// go func() {
	// 	if err := uh.Mailer.SendEmail(); err != nil {
	// 		fmt.Println("unable to send the email:", err)
	//
	// 		return
	// 	}
	// }()
	userResponse := models.NewUserResponse(user)
	token, err := uh.AuthMaker.Create(user.Email, user.UserID, uh.AccessTokenDuration)
	response := models.AuthResponse{
		AccessToken: token,
		User:        userResponse,
	}

	if err := respondWithJSON(w, http.StatusCreated, response); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

func (uh *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	request := models.LoginRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {

		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	user, err := uh.Store.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, errors.New("unable to find user"))
		return
	}
	err = auth.ComparePassword(request.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err)
		return
	}
	userResponse := models.NewUserResponse(user)
	token, err := uh.AuthMaker.Create(user.Email, user.UserID, uh.AccessTokenDuration)
	response := models.AuthResponse{
		AccessToken: token,
		User:        userResponse,
	}
	if err := respondWithJSON(w, http.StatusOK, response); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}
