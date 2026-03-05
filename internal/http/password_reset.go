package http

import (
	"DocNebula/internal/repository"
	"DocNebula/internal/utils"
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type ResetHandler struct {
	UserRepo *repository.UserRepo
}

func (h *ResetHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Email string `json:"email"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	user, err := h.UserRepo.GetByEmail(context.Background(), req.Email)
	if err != nil {
		return
	}

	token, _ := utils.GenerateToken(user.ID)

	utils.SendResetEmail(user.Email, token)

	w.WriteHeader(http.StatusOK)
}

func (h *ResetHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	userID, err := utils.VerifyResetToken(req.Token)
	if err != nil {
		http.Error(w, "invalid token", 400)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)

	_, err = h.UserRepo.DB.Exec(
		`UPDATE users SET password_hash=$1 WHERE id=$2`,
		hash,
		userID,
	)

	if err != nil {
		http.Error(w, "reset failed", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
}
