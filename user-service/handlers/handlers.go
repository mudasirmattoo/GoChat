package handlers

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	"user-service/database"
	"user-service/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type APIResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

var templates map[string]*template.Template

// template cache
func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages := []string{
		"index.html",
		"register.html",
		"login.html",
		"dashboard.html",
	}

	//common templates to all
	baseTemplates := []string{
		"frontend/templates/base.html",
		"frontend/templates/header.html",
		"frontend/templates/footer.html",
	}

	// For each page, create a template set with the base templates
	for _, page := range pages {
		name := page
		ts, err := template.New(name).ParseFiles(append(baseTemplates, "frontend/templates/"+name)...)

		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil

}

func init() {
	var err error
	templates, err = CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache:", err)
	}
}

func RenderTemplate(w http.ResponseWriter, templ string, data interface{}) {
	//get requested template from cache
	t, ok := templates[templ+".html"]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		log.Println("Template not found:", templ)
		return
	}

	err := t.Execute(w, data)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Println("Template error:", err)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	result := database.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	user.Password = ""

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 30),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Registration successful",
		"user":    user,
		"token":   tokenString,
	})

	// http.Redirect(w, r, "/login", http.StatusSeeOther)

	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(user)

	/*  client-side redirection
	-----------------------------------------------------------------
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Success:     true,
		Message:     "Registration successful",
		RedirectURL: "/login", // Client will use this URL to redirect
	})
	*/
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

	var user models.User
	result := database.DB.Where("username = ?", credentials.Username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// yai jwt token generate karega
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // return a  *jwt.Token
		"userID": user.ID,
		"exp":    time.Now().Add(time.Minute * 30).Unix(),
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 30),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Sign in  successful",
		"user":    user,
		"token":   tokenString,
	})
}

func ServeHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/templates/index.html")
}

// func LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var credentials struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 	}

// 	err := json.NewDecoder(r.Body).Decode(&credentials)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	var user models.User
// 	result := database.DB.Where("username = ?", credentials.Username).First(&user)
// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
// 		return
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
// 		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
// 		return
// 	}

// 	http.SetCookie(w, &http.Cookie{
// 		Name:  "flash_message",
// 		Value: "Login successful",
// 		Path:  "/",
// 	})
// 	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
// }
