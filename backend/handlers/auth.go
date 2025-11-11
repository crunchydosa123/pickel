package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"pickel-backend/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	json.NewDecoder(r.Body).Decode(&req)

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	utils.ConnectSupabase()
	db := utils.GetDB()

	var exists bool
	err := db.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", req.Email,
	).Scan(&exists)

	if exists {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	if err != nil {
		fmt.Println("error", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(context.Background(),
		"INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)",
		req.Name, req.Email, string(hashed),
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	json.NewDecoder(r.Body).Decode(&req)

	utils.ConnectSupabase()
	db := utils.GetDB()
	var storedHash string
	var id uuid.UUID
	err := db.QueryRow(context.Background(),
		"SELECT id, password_hash FROM users WHERE email=$1", req.Email,
	).Scan(&id, &storedHash)
	if err != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)) != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	token, err := utils.GenerateJWT(id.String())
	if err != nil {
		http.Error(w, "Token generation failed", 500)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true, // secure from JS access
		Path:     "/",
		Secure:   false, // set true in production (HTTPS)
		SameSite: http.SameSiteNoneMode,
	})

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
	})
}
