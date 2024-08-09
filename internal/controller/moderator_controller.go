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

type ModeratorController interface {
	SendModeratorRequest(w http.ResponseWriter, r *http.Request)
	RemoveModerator(w http.ResponseWriter, r *http.Request)
}

type ModeratorControllerImpl struct {
	ModeratorRepository    repository.ModeratorRepository
	NotificationRepository repository.NotificationRepository
}

func NewModeratorController(moderatorRepo repository.ModeratorRepository, notificationRepo repository.NotificationRepository) ModeratorController {
	return &ModeratorControllerImpl{
		ModeratorRepository:    moderatorRepo,
		NotificationRepository: notificationRepo,
	}
}

func (c *ModeratorControllerImpl) SendModeratorRequest(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		CommunityID       int `json:"community_id"`
		CommunityMemberID int `json:"communitymember_id"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	notification := model.Notification{
		UserID:  requestBody.CommunityMemberID,
		Type:    "request",
		Message: "You have been requested to be a moderator in community with ID " + strconv.Itoa(requestBody.CommunityID),
	}

	if err := c.NotificationRepository.CreateNotification(context.Background(), &notification); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating notification"})
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Moderator request has been sent",
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *ModeratorControllerImpl) RemoveModerator(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moderatorIDstr := vars["id"]

	moderatorID, err := strconv.Atoi(moderatorIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if err := c.ModeratorRepository.DeleteModerator(context.Background(), moderatorID); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error removing moderator"})
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Moderator has been removed",
	}
	httputil.WriteResponse(w, http.StatusOK, response)
}
