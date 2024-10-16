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

type LocationController interface {
	AddLocation(w http.ResponseWriter, r *http.Request)
	UpdateLocation(w http.ResponseWriter, r *http.Request)
	GetLocations(w http.ResponseWriter, r *http.Request)
}

type LocationControllerImpl struct {
	LocationRepository repository.LocationRepository
}

func NewLocationController(locationRepo repository.LocationRepository) LocationController {
	return &LocationControllerImpl{
		LocationRepository: locationRepo,
	}
}

func (c *LocationControllerImpl) AddLocation(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name string `json:"name"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	newLocation := model.Location{
		Name: requestBody.Name,
	}

	if err := c.LocationRepository.AddLocation(context.Background(), &newLocation); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating new location"})
		return
	}

	response := struct {
		Message string         `json:"message"`
		Data    model.Location `json:"data"`
	}{
		Message: "Location has been added",
		Data:    newLocation,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *LocationControllerImpl) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	locationIDstr := vars["id"]

	locationID, err := strconv.Atoi(locationIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid location id"})
		return
	}

	location, err := c.LocationRepository.GetLocationById(context.Background(), locationID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusNotFound, map[string]string{"error": "Location not found"})
		return
	}

	response := struct {
		Message string         `json:"message"`
		Data    model.Location `json:"data"`
	}{
		Message: "Post has been updated",
		Data:    *location,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *LocationControllerImpl) GetLocations(w http.ResponseWriter, r *http.Request) {

	locations, err := c.LocationRepository.GetLocations(context.Background())
	if err != nil {
		httputil.WriteResponse(w, http.StatusOK, map[string]string{"error": "Locations not found"})
		return
	}

	response := struct {
		Message string           `json:"message"`
		Data    []model.Location `json:"data"`
	}{
		Message: "Data has been retrieved successfully",
		Data:    locations,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}
