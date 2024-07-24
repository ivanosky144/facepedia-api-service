package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var allowedExtensions = map[string]bool{
	".jpg": true,
	".png": true,
	".mp4": true,
	".mkv": true,
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	if !allowedExtensions[ext] {
		http.Error(w, "Only images or videos are allowed", http.StatusBadRequest)
		return
	}

	tempFile, err := os.Create(filepath.Join("uploads", handler.Filename))
	if err != nil {
		http.Error(w, "Could not create file", http.StatusBadRequest)
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "File has been uploaded",
	}

	respondJSON(w, http.StatusOK, response)
}
