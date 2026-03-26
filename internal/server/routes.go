package server

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *Handler) {
	r.GET("/health", h.Health)

	api := r.Group("/api")
	{
		tasks := api.Group("/tasks")
		{
			tasks.GET("/", h.GetTasks)
			tasks.POST("/", h.CreateTask)
			tasks.GET("/:id", h.GetTask)
			tasks.DELETE("/:id", h.DeleteTask)
			tasks.GET("/stream", h.StreamTasks)
		}

		export := api.Group("/export")
		{
			export.GET("/", h.ExportGames)
		}

		archive := api.Group("/archive")
		{
			archive.GET("/", h.GetArchiveGames)
			archive.POST("/upload", h.UploadArchiveGames)
		}
	}
}
