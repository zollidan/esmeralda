package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zollidan/esmeralda/models"
	"gorm.io/gorm"
)

func (h *Handlers) FilesRoutes(r chi.Router) {
	r.Get("/", h.GetFiles)
	r.Post("/", h.CreateFile)
	// r.Delete("/{id}", h.DeleteFile)
}

func (h *Handlers) GetFiles(w http.ResponseWriter, r *http.Request) {
	files, err := gorm.G[models.Files](h.DB).Find(context.Background())
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error reading files.",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func (h *Handlers) CreateFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	// Get file from request
	file, fileHandler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content into memory
	// Note: For very large files, consider streaming directly to S3
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read file: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate unique file key (you might want to add timestamp or UUID)
	fileKey := fileHandler.Filename // fmt.Sprintf("%s", fileHandler.Filename)

	// Upload file to S3
	if err := h.S3Client.UploadItem(ctx, fileKey, fileBytes); err != nil {
		http.Error(w, fmt.Sprintf("failed to upload file: %v", err), http.StatusInternalServerError)
		return
	}

	// Construct file URL
	fileURL := fmt.Sprintf("%s/%s/%s", h.Cfg.S3URL, h.S3Client.Bucket, fileKey)

	// Save file metadata to database
	fileRecord := &models.Files{
		Filename: fileHandler.Filename,
		FileURL:  fileURL,
		S3Key:    fileKey,
		FileSize: int64(len(fileBytes)),
	}
	if err := gorm.G[models.Files](h.DB).Create(ctx, fileRecord); err != nil {
		http.Error(w, fmt.Sprintf("failed to save file metadata: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "file uploaded successfully",
		"filename": fileHandler.Filename,
		"url":      fileURL,
		"size":     fileRecord.FileSize,
		"id":       fileRecord.ID,
	})
}

// func (h *Handlers) DeleteFile(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	// Get file ID from URL
// 	idStr := chi.URLParam(r, "id")
// 	id, err := strconv.ParseUint(idStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "invalid file id", http.StatusBadRequest)
// 		return
// 	}

// 	// Get file from database
// 	file, err := gorm.G[models.Files](h.DB).First(ctx, uint(id))
// 	if err != nil {
// 		http.Error(w, "file not found", http.StatusNotFound)
// 		return
// 	}

// 	// Delete from S3 using stored S3Key
// 	if err := h.S3Client.DeleteItem(ctx, file.S3Key); err != nil {
// 		http.Error(w, fmt.Sprintf("failed to delete file from storage: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Delete from database
// 	if err := gorm.G[models.Files](h.DB).Delete(ctx, file); err != nil {
// 		http.Error(w, fmt.Sprintf("failed to delete file record: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message": "file deleted successfully",
// 	})
// }
