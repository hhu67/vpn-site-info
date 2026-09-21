package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type CreatePasswordRequest struct {
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	PasswordConfirm string `json:"password_confirm"`
}

func (h *Handler) CheckPasswordExists(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] CheckPasswordExists called")

	if r.Method != http.MethodGet {
		log.Printf("[HANDLER] CheckPasswordExists: method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var count int
	err := h.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM auth").Scan(&count)
	if err != nil {
		log.Printf("[HANDLER] CheckPasswordExists: database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("[HANDLER] CheckPasswordExists: count=%d", count)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"exists": count > 0})
}

func (h *Handler) CreatePassword(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] CreatePassword called")

	if r.Method != http.MethodPost {
		log.Printf("[HANDLER] CreatePassword: method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var count int
	err := h.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM auth").Scan(&count)
	if err != nil {
		log.Printf("[HANDLER] CreatePassword: database error checking existing password: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		log.Printf("[HANDLER] CreatePassword: password already exists")
		http.Error(w, "Password already exists", http.StatusBadRequest)
		return
	}

	var req CreatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[HANDLER] CreatePassword: invalid request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Password == "" || req.Password != req.PasswordConfirm {
		log.Printf("[HANDLER] CreatePassword: passwords do not match or empty")
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	log.Printf("[HANDLER] CreatePassword: hashing password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[HANDLER] CreatePassword: error hashing password: %v", err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	log.Printf("[HANDLER] CreatePassword: inserting password to database")
	_, err = h.db.Exec(context.Background(), "INSERT INTO auth (password_hash) VALUES ($1)", string(hashedPassword))
	if err != nil {
		log.Printf("[HANDLER] CreatePassword: error saving password: %v", err)
		http.Error(w, "Error saving password", http.StatusInternalServerError)
		return
	}

	log.Printf("[HANDLER] CreatePassword: generating JWT token")
	token, err := h.generateJWT()
	if err != nil {
		log.Printf("[HANDLER] CreatePassword: error generating token: %v", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	log.Printf("[HANDLER] CreatePassword: setting JWT cookie")
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	log.Printf("[HANDLER] CreatePassword: success")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] Login called")

	if r.Method != http.MethodPost {
		log.Printf("[HANDLER] Login: method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[HANDLER] Login: invalid request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("[HANDLER] Login: fetching password hash from database")
	var hashedPassword string
	err := h.db.QueryRow(context.Background(), "SELECT password_hash FROM auth LIMIT 1").Scan(&hashedPassword)
	if err != nil {
		log.Printf("[HANDLER] Login: database error fetching password: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	log.Printf("[HANDLER] Login: comparing passwords")
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		log.Printf("[HANDLER] Login: password mismatch: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	log.Printf("[HANDLER] Login: generating JWT token")
	token, err := h.generateJWT()
	if err != nil {
		log.Printf("[HANDLER] Login: error generating token: %v", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	log.Printf("[HANDLER] Login: setting JWT cookie")
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	log.Printf("[HANDLER] Login: success")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" || req.NewPassword != req.PasswordConfirm {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request or passwords do not match"})
		return
	}

	var hashedPassword string
	err := h.db.QueryRow(context.Background(), "SELECT password_hash FROM auth LIMIT 1").Scan(&hashedPassword)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.OldPassword)); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Incorrect old password"})
		return
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error hashing password"})
		return
	}

	_, err = h.db.Exec(context.Background(), "UPDATE auth SET password_hash = $1", string(newHashedPassword))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error updating password"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handler) generateJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString([]byte(h.jwtSecret))
}
