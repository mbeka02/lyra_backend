package handler

import (
	"errors"
	"net/http"

	"github.com/mbeka02/lyra_backend/internal/model"
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

func (h *DoctorHandler) HandleListMyPatients(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	patients, err := h.doctorService.ListPatientsUnderCare(r.Context(), payload.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, patients)
}

func (h *DoctorHandler) HandleGetDoctors(w http.ResponseWriter, r *http.Request) {
	params := NewQueryParamExtractor(r)
	page := params.GetInt32("page", 0)
	pageSize := int32(10)
	offset := page * pageSize

	response, err := h.doctorService.GetDoctors(r.Context(), params.GetString("county"), params.GetString("specialization"), params.GetString("minPrice"), params.GetString("maxPrice"), params.GetString("sort"), params.GetString("order"), params.GetInt32("minExperience", 0), params.GetInt32("maxExperience", 10000), pageSize, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("unable to get doctor details"))
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *DoctorHandler) HandleCreateDoctor(w http.ResponseWriter, r *http.Request) {
	request := model.CreateDoctorRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	response, err := h.doctorService.CreateDoctor(r.Context(), request, payload.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusCreated, response)
}
