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
	AddReview(w http.ResponseWriter, r *http.Request)
	GetUniversityReviews(w http.ResponseWriter, r *http.Request)
}

type UniversityControllerImpl struct {
	UniversityRepository repository.UniversityRepository
	ReviewRepository     repository.ReviewRepository
}

func NewUniversityController(universityRepo repository.UniversityRepository, reviewRepo repository.ReviewRepository) UniversityController {
	return &UniversityControllerImpl{
		UniversityRepository: universityRepo,
		ReviewRepository:     reviewRepo,
	}
}

func (c *UniversityControllerImpl) AddUniversity(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name       string `json:"name"`
		Summary    string `json:"summary"`
		LocationID int    `json:"location_id"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	newUniversity := model.University{
		Name:       requestBody.Name,
		Summary:    requestBody.Summary,
		LocationID: requestBody.LocationID,
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

	universities, err := c.UniversityRepository.GetUniversities(context.Background())
	if err != nil {
		httputil.WriteResponse(w, http.StatusOK, map[string]string{"error": "Universities not found"})
		return
	}

	response := struct {
		Message string             `json:"message"`
		Data    []model.University `json:"data"`
	}{
		Message: "Data has been retrieved successfully",
		Data:    universities,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *UniversityControllerImpl) AddReview(w http.ResponseWriter, r *http.Request) {

	var requestBody struct {
		UserID       int    `json:"user_id"`
		UniversityID int    `json:"university_id"`
		Text         string `json:"text"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	newUniversityReview := model.Review{
		UserID:       requestBody.UserID,
		UniversityID: requestBody.UniversityID,
		Text:         requestBody.Text,
	}

	if err := c.ReviewRepository.CreateReview(context.Background(), &newUniversityReview); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating new review"})
		return
	}

	response := struct {
		Message string       `json:"message"`
		Data    model.Review `json:"data"`
	}{
		Message: "New review has been added",
		Data:    newUniversityReview,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *UniversityControllerImpl) GetUniversityReviews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	universityIDstr := vars["university_id"]

	universityID, err := strconv.Atoi(universityIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid university id"})
		return
	}

	universityReviews, err := c.ReviewRepository.GetReviewsByUniversityID(context.Background(), universityID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving university reviews"})
		return
	}

	response := struct {
		Message string         `json:"message"`
		Data    []model.Review `json:"data"`
	}{
		Message: "Data has been retrieved successfully",
		Data:    universityReviews,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}
