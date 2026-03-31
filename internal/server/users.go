package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zollidan/esmeralda/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"` //nolint:gosec // Payload field
}

type LoginUserResponse struct {
	AccessToken  string `json:"access_token"`  //nolint:gosec // Payload field
	RefreshToken string `json:"refresh_token"` //nolint:gosec // Payload field
}

// PostLoginUser godoc
// @Summary      Авторизация пользователя
// @Description  Проверяет учетные данные и возвращает JWT токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginUserRequest   true  "Логин и пароль"
// @Success      200      {object}  LoginUserResponse
// @Failure      400      {object}  errorResponse
// @Failure      401      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /auth/login [post]
func (h *Handler) PostLoginUser(c *gin.Context) {
	var req LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.users.FindByUsername(c.Request.Context(), req.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	if user == nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"exp":        time.Now().Add(time.Duration(h.cfg.Auth.AccessTokenTTL) * time.Second).Unix(),
		"token_type": "access",
	})

	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"exp":        time.Now().Add(time.Duration(h.cfg.Auth.RefreshTokenTTL) * time.Second).Unix(),
		"token_type": "refresh",
	})

	access_token_str, err := access_token.SignedString([]byte(h.cfg.Auth.JWTSecret))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	refresh_token_str, err := refresh_token.SignedString([]byte(h.cfg.Auth.JWTSecret))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate refresh token"})
		return
	}

	if err := h.refreshTokens.Create(c.Request.Context(), &models.RefreshToken{
		UserID:    user.ID,
		Token:     refresh_token_str,
		ExpiresAt: time.Now().Add(time.Duration(h.cfg.Auth.RefreshTokenTTL) * time.Second),
	}); err != nil {
		c.JSON(500, gin.H{"error": "failed to save refresh token"})
		return
	}

	c.JSON(200, LoginUserResponse{
		AccessToken:  access_token_str,
		RefreshToken: refresh_token_str,
	})
}
