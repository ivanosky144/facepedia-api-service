package controller

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	httputil "github.com/temuka-api-service/pkg/http"
)

type FileUploadHandler struct {
	AllowedExtensions map[string]bool
	UploadDirectory   string
}

func NewFileUploadController(uploadDir string) *FileUploadHandler {
	return &FileUploadHandler{
		AllowedExtensions: map[string]bool{
			".jpg": true,
			".png": true,
			".mp4": true,
			".mkv": true,
		},
		UploadDirectory: uploadDir,
	}
}

func (h *FileUploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Could not parse multipart form"})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Could not get uploaded file"})
		return
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	if !h.AllowedExtensions[ext] {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Only images or videos are allowed"})
		return
	}

	tempFilePath := filepath.Join(h.UploadDirectory, handler.Filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not create file"})
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not save file"})
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "File has been uploaded",
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}
