package controller

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/temuka-api-service/config"
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

	s3Key := fmt.Sprintf("uploads/%s", handler.Filename)

	_, err = config.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(config.S3Bucket),
		Key:         aws.String(s3Key),
		Body:        file,
		ContentType: aws.String(handler.Header.Get("Content-Type")),
		ACL:         types.ObjectCannedACLPublicRead,
	})

	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not upload file to S3"})
		return
	}

	response := struct {
		Message string `json:"message"`
		URL     string `json:"url"`
	}{
		Message: "File has been uploaded",
		URL:     fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", config.S3Bucket, config.S3Client.Options().Region, s3Key),
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}
