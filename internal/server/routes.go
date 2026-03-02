package server

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, h *Handler) {
	api := r.Group("/api")
	{
		api.POST("/tasks", h.CreateTask)
		api.GET("/tasks", h.GetTasks)
	}
}