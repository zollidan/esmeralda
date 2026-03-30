package server

import "github.com/gin-gonic/gin"

// GetArchiveGames godoc
// @Summary      Получить архивные матчи
// @Description  Возвращает архивные матчи (эндпоинт в разработке)
// @Tags         archive
// @Produce      json
// @Success      501  {object}  errorResponse
// @Security     BearerAuth
// @Router       /archive/ [get]
func (h *Handler) GetArchiveGames(c *gin.Context) {}

// UploadArchiveGames godoc
// @Summary      Загрузить архивные матчи
// @Description  Загружает архивные матчи из файла (эндпоинт в разработке)
// @Tags         archive
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Файл архива"
// @Success      501   {object}  errorResponse
// @Security     BearerAuth
// @Router       /archive/upload [post]
func (h *Handler) UploadArchiveGames(c *gin.Context) {}
