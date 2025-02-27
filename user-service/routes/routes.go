package routes

import (
	"net/http"
	"user-service/handlers"
)

func RegisterRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		handlers.RenderTemplate(w, "index", map[string]string{"Title": "Home"})
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderTemplate(w, "register", map[string]string{"Title": "Register"})
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderTemplate(w, "login", map[string]string{"Title": "Login"})
	})

	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderTemplate(w, "dashboard", map[string]string{"Title": "Dashboard"})
	})

	http.HandleFunc("/register-user", handlers.RegisterHandler)
	http.HandleFunc("/login-user", handlers.LoginHandler)

}
