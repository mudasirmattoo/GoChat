package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"user-service/database"
	"user-service/models"

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

	user.Password = ""

	http.SetCookie(w, &http.Cookie{
		Name:  "flash_message",
		Value: "Registration Ssuccessful",
		Path:  "/",
	}) //SetCookie(w http.ResponseWriter, cookie *http.Cookie)  -->SetCookie adds a Set-Cookie header to the provided [ResponseWriter]'s header
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

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

// LoginHandler handles user login
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

	// Fetch user from DB using GORM
	var user models.User
	result := database.DB.Where("username = ?", credentials.Username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "flash_message",
		Value: "Login successful",
		Path:  "/",
	})
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

	fmt.Fprintln(w, "Login successful")
}

// ServeHomePage serves the homepage
func ServeHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/templates/index.html")
}
