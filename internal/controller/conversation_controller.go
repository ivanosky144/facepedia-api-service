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

type ConversationController interface {
	AddConversation(w http.ResponseWriter, r *http.Request)
	DeleteConversation(w http.ResponseWriter, r *http.Request)
	GetConversationsByUserID(w http.ResponseWriter, r *http.Request)
	AddParticipant(w http.ResponseWriter, r *http.Request)
	AddMessage(w http.ResponseWriter, r *http.Request)
}

type ConversationControllerImpl struct {
	ConversationRepository repository.ConversationRepository
	UserRepository         repository.UserRepository
}

func NewConversationController(conversationRepo repository.ConversationRepository, userRepo repository.UserRepository) ConversationController {
	return &ConversationControllerImpl{
		ConversationRepository: conversationRepo,
		UserRepository:         userRepo,
	}
}

func (c *ConversationControllerImpl) AddConversation(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Title  string `json:"title"`
		UserID int    `json:"user_id"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	newConversation := model.Conversation{
		UserID: requestBody.UserID,
		Title:  requestBody.Title,
	}

	if err := c.ConversationRepository.CreateConversation(context.Background(), &newConversation); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating conversation"})
		return
	}

	response := struct {
		Message string             `json:"message"`
		Data    model.Conversation `json:"data"`
	}{
		Message: "Conversation has been created",
		Data:    newConversation,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *ConversationControllerImpl) AddMessage(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		ParticipantID int    `json:"participant_id"`
		Text          string `json:"text"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	newMessage := model.Message{
		ParticipantID: requestBody.ParticipantID,
		Text:          requestBody.Text,
	}

	if err := c.ConversationRepository.AddMessage(context.Background(), &newMessage); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error adding message"})
		return
	}

	response := struct {
		Message string        `json:"message"`
		Data    model.Message `json:"data"`
	}{
		Message: "Message has been created",
		Data:    newMessage,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *ConversationControllerImpl) GetConversationsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstr := vars["user_id"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user id"})
		return
	}

	conversations, err := c.ConversationRepository.GetConversationsByUserID(context.Background(), userID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving conversatiosn"})
		return
	}

	response := struct {
		Message string               `json:"message"`
		Data    []model.Conversation `json:"data"`
	}{
		Message: "Conversations have been retrieved",
		Data:    conversations,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *ConversationControllerImpl) AddParticipant(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		ConversationID int `json:"conversation_id"`
		UserID         int `json:"user_id"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	newParticipant := model.Participant{
		UserID:         requestBody.UserID,
		ConversationID: requestBody.ConversationID,
	}

	if err := c.ConversationRepository.AddParticipant(context.Background(), &newParticipant); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error adding participant"})
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Participant has been added",
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *ConversationControllerImpl) DeleteConversation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationIDstr := vars["id"]

	conversationID, err := strconv.Atoi(conversationIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid conversation id"})
		return
	}

	if err := c.ConversationRepository.DeleteConversation(context.Background(), conversationID); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error deleting conversation"})
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Conversation has been deleted",
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}
