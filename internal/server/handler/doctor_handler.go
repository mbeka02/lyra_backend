package handler

import (
	"errors"
	"net/http"

	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/middleware"
	"github.com/mbeka02/lyra_backend/internal/server/service"
)

type DoctorHandler struct {
	doctorService service.DoctorService
}

func NewDoctorHandler(doctorService service.DoctorService) *DoctorHandler {
	return &DoctorHandler{
		doctorService,
	}
}

func (h *DoctorHandler) HandleGetDoctors(w http.ResponseWriter, r *http.Request) {
	params := NewQueryParamExtractor(r)
	page := params.GetInt("page", 0)
	pageSize := int32(10)
	offset := int32(page) * pageSize

	response, err := h.doctorService.GetDoctors(r.Context(), params.GetString("county"), params.GetString("specialization"), params.GetString("minPrice"), params.GetString("maxPrice"), params.GetString("sort"), params.GetString("order"), params.GetInt32("minExperience", 0), params.GetInt32("maxExperience", 10000), pageSize, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("unable to get doctor details"))
	}

	if err := respondWithJSON(w, http.StatusOK, response); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
	}
}

func (h *DoctorHandler) HandleCreateDoctor(w http.ResponseWriter, r *http.Request) {
	request := model.CreateDoctorRequest{}
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
	response, err := h.doctorService.CreateDoctor(r.Context(), request, payload.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	if err := respondWithJSON(w, http.StatusCreated, response); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}
