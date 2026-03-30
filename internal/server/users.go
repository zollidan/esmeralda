package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

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

	// проверка пароля
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(h.cfg.Auth.JWTSecret))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, LoginUserResponse{
		Token: ss,
	})
}