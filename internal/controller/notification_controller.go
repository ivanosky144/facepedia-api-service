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

type NotificationController interface {
	GetNotificationsByUser(w http.ResponseWriter, r *http.Request)
}

type NotificationControllerImpl struct {
	NotificationRepository repository.NotificationRepository
}

func NewNotificationController(notifRepo repository.NotificationRepository) NotificationController {
	return &NotificationControllerImpl{
		NotificationRepository: notifRepo,
	}
}

func (c *NotificationControllerImpl) GetNotificationsByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstr := vars["user_id"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user id"})
		return
	}

	notifications, err := c.NotificationRepository.GetNotificationsByUserID(context.Background(), userID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving notifications"})
		return
	}

	response := struct {
		Message string               `json:"message"`
		Data    []model.Notification `json:"data"`
	}{
		Message: "User posts have been retrieved",
		Data:    notifications,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}
