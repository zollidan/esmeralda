package server

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *Handler) {
	r.GET("/health", h.Health)

	api := r.Group("/api")
	{
		api.POST("/tasks", h.CreateTask)
		api.GET("/tasks", h.GetTasks)
		api.GET("/tasks/:id", h.GetTask)
		api.GET("/export", h.ExportGames)
		api.GET("/progress/ws", h.WSHandler)
		
		archive := api.Group("/archive")
		{
			archive.GET("/", h.GetArchiveGames)
			archive.POST("/upload", h.UploadArchiveGames)
		}
	}
}
