package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zollidan/esmeralda/models"
	"github.com/zollidan/esmeralda/schemas"
	"github.com/zollidan/esmeralda/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (h *Handlers) AuthRoutes(r chi.Router) {
	r.Post("/login", h.LoginUser)
	r.Post("/register", h.RegisterUser)
}

func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req *schemas.LoginUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error invalid JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, map[string]string{"message": "login endpoint"})

}

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req *schemas.RegisterUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error invalid JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	_, err = gorm.G[models.User](h.DB).Where("username = ?", req.Username).First(context.Background())
	if err == nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, fmt.Sprintf("Error checking existing user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error hashing password: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	err = gorm.G[models.User](h.DB).Create(context.Background(), user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user created"})
}
