package http

import (
	"DocNebula/internal/repository"
	"DocNebula/internal/utils"
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo *repository.UserRepo
}

type authReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {

	var req authReq
	json.NewDecoder(r.Body).Decode(&req)

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)

	user, err := h.UserRepo.Create(context.Background(), req.Email, string(hash))
	if err != nil {
		http.Error(w, "signup failed", 500)
		return
	}

	token, _ := utils.GenerateToken(user.ID)

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req authReq
	json.NewDecoder(r.Body).Decode(&req)

	user, err := h.UserRepo.GetByEmail(context.Background(), req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", 401)
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		http.Error(w, "invalid credentials", 401)
		return
	}

	token, _ := utils.GenerateToken(user.ID)

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
