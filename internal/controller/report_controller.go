package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/repository"
	httputil "github.com/temuka-api-service/pkg/http"
)

type ReportController interface {
	CreateReport(w http.ResponseWriter, r *http.Request)
	DeleteReport(w http.ResponseWriter, r *http.Request)
}

type ReportControllerImpl struct {
	ReportRepository repository.ReportRepository
}

func NewReportController(reportRepo repository.ReportRepository) ReportController {
	return &ReportControllerImpl{
		ReportRepository: reportRepo,
	}
}

func (c *ReportControllerImpl) CreateReport(w http.ResponseWriter, r *http.Request) {

}

func (c *ReportControllerImpl) DeleteReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportIDstr := vars["id"]

	reportID, err := strconv.Atoi(reportIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid report id"})
		return
	}

	if err := c.ReportRepository.DeleteReport(context.Background(), reportID); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error deleting comment"})
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Report has been deleted",
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}
