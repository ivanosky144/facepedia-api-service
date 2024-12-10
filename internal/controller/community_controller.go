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

type CommunityController interface {
	CreateCommunity(w http.ResponseWriter, r *http.Request)
	GetCommunities(w http.ResponseWriter, r *http.Request)
	DeleteCommunity(w http.ResponseWriter, r *http.Request)
	UpdateCommunity(w http.ResponseWriter, r *http.Request)
	GetUserJoinedCommunities(w http.ResponseWriter, r *http.Request)
	JoinCommunity(w http.ResponseWriter, r *http.Request)
	GetCommunityPosts(w http.ResponseWriter, r *http.Request)
	GetCommunityDetail(w http.ResponseWriter, r *http.Request)
}

type CommunityControllerImpl struct {
	CommunityRepository repository.CommunityRepository
}

func NewCommunityController(repo repository.CommunityRepository) CommunityController {
	return &CommunityControllerImpl{
		CommunityRepository: repo,
	}
}

func (c *CommunityControllerImpl) CreateCommunity(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		LogoPicture string `json:"logo_picture"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if !c.CommunityRepository.CheckCommunityNameAvailability(context.Background(), requestBody.Name) {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Community with the same name already exist"})
		return
	}

	newCommunity := model.Community{
		Name:        requestBody.Name,
		Description: requestBody.Description,
		LogoPicture: requestBody.LogoPicture,
	}

	if err := c.CommunityRepository.CreateCommunity(context.Background(), &newCommunity); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating community"})
		return
	}

	response := struct {
		Message string          `json:"message"`
		Data    model.Community `json:"data"`
	}{
		Message: "Community has been created",
		Data:    newCommunity,
	}
	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *CommunityControllerImpl) GetCommunities(w http.ResponseWriter, r *http.Request) {

	communities, err := c.CommunityRepository.GetCommunities(context.Background())
	if err != nil {
		httputil.WriteResponse(w, http.StatusNotFound, map[string]string{"error": "Error retrieving communities"})
		return
	}

	response := struct {
		Message string            `json:"message"`
		Data    []model.Community `json:"data"`
	}{
		Message: "Communities have been retireved",
		Data:    communities,
	}
	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *CommunityControllerImpl) UpdateCommunity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityIDstr := vars["id"]
	communityID, err := strconv.Atoi(communityIDstr)

	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Invalid community id"})
		return
	}

	var requestBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		LogoPicture string `json:"logo_picture"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	updatedCommunity := model.Community{
		Name:        requestBody.Name,
		Description: requestBody.Description,
		LogoPicture: requestBody.LogoPicture,
	}

	if err := c.CommunityRepository.UpdateCommunity(context.Background(), communityID, &updatedCommunity); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error updating community"})
		return
	}

	response := struct {
		Message string          `json:"message"`
		Data    model.Community `json:"data"`
	}{
		Message: "Communities have been retireved",
		Data:    updatedCommunity,
	}
	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *CommunityControllerImpl) JoinCommunity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityIDstr := vars["community_id"]

	communityID, err := strconv.Atoi(communityIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid community id"})
		return
	}

	var requestBody struct {
		UserID int `json:"user_id"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	community, err := c.CommunityRepository.GetCommunityDetailByID(context.Background(), communityID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving community"})
		return
	}
	if community == nil {
		httputil.WriteResponse(w, http.StatusNotFound, map[string]string{"error": "Community not found"})
		return
	}

	existingMember, err := c.CommunityRepository.CheckMembership(context.Background(), communityID, requestBody.UserID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error checking community membership"})
		return
	}
	if existingMember != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "User already a member of the community"})
		return
	}

	newMember := model.CommunityMember{
		UserID:      requestBody.UserID,
		CommunityID: communityID,
	}

	if err := c.CommunityRepository.AddCommunityMember(context.Background(), &newMember); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error adding community member"})
		return
	}

	community.MembersCount++
	if err := c.CommunityRepository.UpdateCommunity(context.Background(), communityID, community); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error updating community"})
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Successfully joined the community",
	}
	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *CommunityControllerImpl) DeleteCommunity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityIDstr := vars["id"]

	communityID, err := strconv.Atoi(communityIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid community"})
		return
	}

	if err := c.CommunityRepository.DeleteCommunity(context.Background(), communityID); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error deleting community"})
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Community has been deleted",
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *CommunityControllerImpl) GetCommunityPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityIDstr := vars["id"]

	communityID, err := strconv.Atoi(communityIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid community"})
		return
	}

	filters := make(map[string]interface{})
	if topic := r.URL.Query().Get("topic"); topic != "" {
		filters["topic"] = topic
	}
	if sort := r.URL.Query().Get("sort"); sort != "" {
		filters["sort"] = sort
	}
	if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
		filters["sort_by"] = sortBy
	}

	communityPosts, err := c.CommunityRepository.GetCommunityPosts(context.Background(), communityID, filters)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving posts"})
		return
	}

	response := struct {
		Message string                `json:"message"`
		Data    []model.CommunityPost `json:"data"`
	}{
		Message: "Community posts has been retrieved",
		Data:    communityPosts,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *CommunityControllerImpl) GetCommunityDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityIDstr := vars["id"]

	communityID, err := strconv.Atoi(communityIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid community id"})
		return
	}

	community, err := c.CommunityRepository.GetCommunityDetailByID(context.Background(), communityID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving community detail"})
		return
	}

	response := struct {
		Message string          `json:"message"`
		Data    model.Community `json:"data"`
	}{
		Message: "Community detail has been retrieved",
		Data:    *community,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *CommunityControllerImpl) GetUserJoinedCommunities(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		UserID int `json:"user_id"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	userCommunities, err := c.CommunityRepository.GetUserJoinedCommunities(context.Background(), requestBody.UserID)

	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving communities"})
		return
	}

	response := struct {
		Message string            `json:"message"`
		Data    []model.Community `json:"data"`
	}{
		Message: "Communities has been retrieved",
		Data:    userCommunities,
	}
	httputil.WriteResponse(w, http.StatusOK, response)
}
