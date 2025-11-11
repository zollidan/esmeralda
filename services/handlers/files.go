package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zollidan/esmeralda/database"
	"github.com/zollidan/esmeralda/models"
	"gorm.io/gorm"
)

func (h *Handlers) FilesRoutes(r chi.Router) {
	r.Get("/", h.GetFiles)
	r.Post("/", h.CreateFile)
	r.Get("/{id}", h.GetFile)
	r.Get("/{id}/download", h.DownloadFile)
	r.Delete("/{id}", h.DeleteFile)
}

func (h *Handlers) GetFiles(w http.ResponseWriter, r *http.Request) {
	database.Get[models.Files](w, r, h.DB)
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

func (h *Handlers) GetFile(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "id")
	database.GetByID[models.Files](w, r, h.DB, fileID)
}

func (h *Handlers) DownloadFile(w http.ResponseWriter, r *http.Request) {
// 	fileID := chi.URLParam(r, "id")

// 	result, err := gorm.G[models.Files](h.DB).Where("id = ?", fileID).First(context.Background())
// 	if err != nil {
// 		http.Error(w, "File not found", http.StatusNotFound)
// 		return
// 	}

// 	// Stream file content to response
// 	obj, err := h.S3Client.GetItem(context.Background(), result.S3Key)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("failed to get file from storage: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	defer obj.Body.Close()

// 	if obj.ContentType != nil {
// 		w.Header().Set("Content-Type", *obj.ContentType)
// 	} else {
// 		w.Header().Set("Content-Type", "application/octet-stream")
// 	}

// 	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", result.Filename))

// 	if obj.ContentLength != nil {
// 		w.Header().Set("Content-Length", fmt.Sprintf("%d", *obj.ContentLength))
// 	}

// 	_, err = io.Copy(w, obj.Body)
// 	if err != nil {
// 		log.Printf("Error streaming file: %v\n", err)
// 		return
// 	}
}

func (h *Handlers) DeleteFile(w http.ResponseWriter, r *http.Request) {
	// fileID := chi.URLParam(r, "id")

	// result, err := gorm.G[models.Files](h.DB).Where("id = ?", fileID).First(context.Background())
	// if err != nil {
	// 	http.Error(w, "File not found", http.StatusNotFound)
	// 	return
	// }

	// // Delete file from S3
	// if err := h.S3Client.DeleteItem(context.Background(), result.S3Key); err != nil {
	// 	http.Error(w, fmt.Sprintf("failed to delete file from storage: %v", err), http.StatusInternalServerError)
	// 	return
	// }

	// // Delete file record from database
	// _, err = gorm.G[models.Files](h.DB).Where("id = ?", fileID).Delete(context.Background())
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("failed to delete file record: %v", err), http.StatusInternalServerError)
	// 	return
	// }

	// utils.ResponseJSON(w, http.StatusNoContent, map[string]string{
	// 	"message": "file deleted successfully",
	// })
}
