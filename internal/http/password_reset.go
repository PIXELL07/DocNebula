package http

import (
	"DocNebula/internal/repository"
	"DocNebula/internal/utils"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type ResetHandler struct {
	UserRepo *repository.UserRepo
}

type forgotReq struct {
	Email string `json:"email"`
}

type resetReq struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (h *ResetHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {

	var req forgotReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "email required", http.StatusBadRequest)
		return
	}

	user, err := h.UserRepo.GetByEmail(r.Context(), req.Email)

	if err != nil || user == nil {
		// don't reveal if email exists
		w.WriteHeader(http.StatusOK)
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	err = utils.SendResetEmail(user.Email, token)
	if err != nil {
		http.Error(w, "email send failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "reset email sent",
	})
}

func (h *ResetHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {

	var req resetReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" || req.Password == "" {
		http.Error(w, "token and password required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, "password too short", http.StatusBadRequest)
		return
	}

	userID, err := utils.VerifyResetToken(req.Token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "password hashing failed", http.StatusInternalServerError)
		return
	}

	_, err = h.UserRepo.DB.ExecContext(
		r.Context(),
		`UPDATE users SET password_hash=$1 WHERE id=$2`,
		hash,
		userID,
	)

	if err != nil {
		http.Error(w, "reset failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "password updated",
	})
}
