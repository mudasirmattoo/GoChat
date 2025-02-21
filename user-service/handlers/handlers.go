package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"user-service/database"
	"user-service/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}
	// Temporary struct to handle input validation
	// var input struct {
	// 	ID              int    `json:"id"`
	// 	Username        string `json:"username"`
	// 	Email           string `json:"email"`
	// 	Password        string `json:"password"`
	// 	ConfirmPassword string `json:"confirm_password"`
	// }

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// if input.Password != input.ConfirmPassword {
	// 	http.Error(w, "Passwords do not match", http.StatusBadRequest)
	// 	return
	// }

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Create user in DB using GORM

	result := database.DB.Create(&user)

	if result.Error != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	// err = SaveUser(user)
	// if err != nil {
	// 	http.Error(w, "User already exists", http.StatusConflict)
	// 	return
	// }

	// Remove password before returning response
	user.Password = ""
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
	fmt.Fprintln(w, "User registered successfully")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// user, exists := GetUser(credentials.Username)
	// if !exists || !ComparePassword(user.Password, credentials.Password) {
	// 	http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	// 	return
	// }

	// Fetch user from DB using GORM
	var user models.User
	result := database.DB.Where("username = ?", credentials.Username).First(&user)
	if result.Error != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Login successful")
}

func ServeHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/index.html")

}
