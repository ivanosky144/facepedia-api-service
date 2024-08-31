package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/repository"
	httputil "github.com/temuka-api-service/pkg/http"
)

type UniversityController interface {
	AddUniversity(w http.ResponseWriter, r *http.Request)
	UpdateUniversity(w http.ResponseWriter, r *http.Request)
	GetUniversityDetail(w http.ResponseWriter, r *http.Request)
	GetUniversities(w http.ResponseWriter, r *http.Request)
}

type UniversityControllerImpl struct {
	UniversityRepository repository.UniversityRepository
}

func NewUniversityController(universityRepo repository.UniversityRepository) UniversityController {
	return &UniversityControllerImpl{
		UniversityRepository: universityRepo,
	}
}

func (c *UniversityControllerImpl) AddUniversity(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name     string `json:"name"`
		Summary  string `json:"summary"`
		Location string `json:"location"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	newUniversity := model.University{
		Name:     requestBody.Name,
		Summary:  requestBody.Summary,
		Location: requestBody.Location,
	}

	if err := c.UniversityRepository.CreateUniversity(context.Background(), &newUniversity); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating new university"})
		return
	}

	response := struct {
		Message string           `json:"message"`
		Data    model.University `json:"data"`
	}{
		Message: "University has been added",
		Data:    newUniversity,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *UniversityControllerImpl) UpdateUniversity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	universityIDstr := vars["id"]

	universityID, err := strconv.Atoi(universityIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid university id"})
		return
	}

	university, err := c.UniversityRepository.GetUniversityDetailByID(context.Background(), universityID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusNotFound, map[string]string{"error": "University not found"})
		return
	}

	response := struct {
		Message string           `json:"message"`
		Data    model.University `json:"data"`
	}{
		Message: "Post has been updated",
		Data:    *university,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *UniversityControllerImpl) GetUniversityDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	universityIDstr := vars["id"]

	universityID, err := strconv.Atoi(universityIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid university id"})
		return
	}

	university, err := c.UniversityRepository.GetUniversityDetailByID(context.Background(), universityID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusNotFound, map[string]string{"error": "University not found"})
		return
	}

	response := struct {
		Message string           `json:"message"`
		Data    model.University `json:"data"`
	}{
		Message: "Post has been updated",
		Data:    *university,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *UniversityControllerImpl) GetUniversities(w http.ResponseWriter, r *http.Request) {

}
