package server

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zollidan/esmeralda/internal/static"
)

func SetupRoutes(r *gin.Engine, h *Handler) {
	api := r.Group("/api")
	{
		api.POST("/tasks", h.CreateTask)
		api.GET("/tasks", h.GetTasks)
		api.GET("/export", h.ExportGames)
	}

	subFS, _ := fs.Sub(static.StaticFiles, "dist")
	r.NoRoute(func(ctx *gin.Context) {
		staticServer := http.FileServer(http.FS(subFS))
		staticServer.ServeHTTP(ctx.Writer, ctx.Request)
	})
}
