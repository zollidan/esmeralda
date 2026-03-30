package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/zollidan/esmeralda/docs"
	"github.com/zollidan/esmeralda/internal/middleware"
)

func SetupRoutes(r *gin.Engine, h *Handler) {
	r.GET("/health", h.Health)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/api/auth")
	{
		auth.POST("/login", h.PostLoginUser)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(h.cfg.Auth.JWTSecret))
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
