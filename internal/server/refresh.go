package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zollidan/esmeralda/internal/models"
)

type RefreshRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

func (h *Handler) PostRefreshToken(c *gin.Context) {
    var req RefreshRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }

    token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(h.cfg.Auth.JWTSecret), nil
    })
    if err != nil || !token.Valid {
        c.JSON(401, gin.H{"error": "invalid refresh token"})
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || claims["token_type"] != "refresh" {
        c.JSON(401, gin.H{"error": "invalid token type"})
        return
    }

    stored, err := h.refreshTokens.FindByToken(c.Request.Context(), req.RefreshToken)
    if err != nil {
        c.JSON(500, gin.H{"error": "internal error"})
        return
    }
    if stored == nil {
        c.JSON(401, gin.H{"error": "refresh token not found"})
        return
    }
    if time.Now().After(stored.ExpiresAt) {
        c.JSON(401, gin.H{"error": "refresh token expired"})
        return
    }

    userID := uint(claims["user_id"].(float64))

    if err := h.refreshTokens.DeleteByToken(c.Request.Context(), req.RefreshToken); err != nil {
        c.JSON(500, gin.H{"error": "internal error"})
        return
    }

    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":    userID,
        "exp":        time.Now().Add(time.Duration(h.cfg.Auth.AccessTokenTTL) * time.Second).Unix(),
        "token_type": "access",
    })
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":    userID,
        "exp":        time.Now().Add(time.Duration(h.cfg.Auth.RefreshTokenTTL) * time.Second).Unix(),
        "token_type": "refresh",
    })

    accessTokenStr, err := accessToken.SignedString([]byte(h.cfg.Auth.JWTSecret))
    if err != nil {
        c.JSON(500, gin.H{"error": "failed to generate access token"})
        return
    }
    refreshTokenStr, err := refreshToken.SignedString([]byte(h.cfg.Auth.JWTSecret))
    if err != nil {
        c.JSON(500, gin.H{"error": "failed to generate refresh token"})
        return
    }

    if err := h.refreshTokens.Create(c.Request.Context(), &models.RefreshToken{
        UserID:    userID,
        Token:     refreshTokenStr,
        ExpiresAt: time.Now().Add(time.Duration(h.cfg.Auth.RefreshTokenTTL) * time.Second),
    }); err != nil {
        c.JSON(500, gin.H{"error": "failed to save refresh token"})
        return
    }

    c.JSON(200, RefreshResponse{
        AccessToken:  accessTokenStr,
        RefreshToken: refreshTokenStr,
    })
}