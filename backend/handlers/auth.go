package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"pickel-backend/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	json.NewDecoder(r.Body).Decode(&req)

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	utils.ConnectSupabase()
	db := utils.GetDB()

	_, err := db.Exec(context.Background(),
		"INSERT INTO users (email, password_hash) VALUES ($1, $2)",
		req.Email, string(hashed),
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

	json.NewEncoder(w).Encode(map[string]string{"token": token})

}
