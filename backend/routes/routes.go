package routes

import (
	"backend/auth"
	"backend/handlers"
	"net/http"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		handlers.RenderTemplate(w, "index", map[string]string{"Title": "Home"})
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderTemplate(w, "register", map[string]string{"Title": "Register"})
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderTemplate(w, "login", map[string]string{"Title": "Login"})
	})

	mux.HandleFunc("/logout", handlers.LogoutHandler)

	mux.HandleFunc("/dashboard", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderTemplate(w, "dashboard", map[string]string{"Title": "Dashboard"})
	}))

	mux.HandleFunc("/register-user", handlers.RegisterHandler)
	mux.HandleFunc("/login-user", handlers.LoginHandler)

	return mux
}

// func RegisterRoutes() *mux.Router {
// 	r := mux.NewRouter()

// 	r.HandleFunc("/create/comment", handlers.CreateComment).Methods("POST")
// 	r.HandleFunc("/get/comments", handlers.GetComments).Methods("GET")
// 	r.HandleFunc("/get/comments/{post_id}", handlers.GetCommentsByPostID).Methods("GET")
// 	r.HandleFunc("/delete/comments/{id}", handlers.DeleteComment).Methods("DELETE")
// 	r.HandleFunc("/edit/comments/{id}", handlers.EditComment).Methods("PUT")

// 	return r
// }

// func RegisterRoutes() *mux.Router {
// 	r := mux.NewRouter()

// 	r.HandleFunc("/", handlers.ServeHome).Methods("GET")
// 	r.HandleFunc("/create/post", handlers.CreatePost).Methods("POST")
// 	r.HandleFunc("/get/posts", handlers.GetPost).Methods("GET")
// 	r.HandleFunc("/get/posts/{id}", handlers.GetPostByID).Methods("GET")
// 	r.HandleFunc("/update/posts/{id}", handlers.UpdatePost).Methods("PUT")
// 	r.HandleFunc("/delete/posts/{id}", handlers.DeletePost).Methods("DELETE")

// 	return r
// }

// func RegisterRoutes() *mux.Router {
// 	r := mux.NewRouter()

// 	r.HandleFunc("/like/add", handlers.AddLike).Methods("POST")
// 	r.HandleFunc("/like/{post_id}", handlers.GetLikesByPostID).Methods("GET")
// 	r.HandleFunc("/like/remove/{id}", handlers.RemoveLike).Methods("DELETE")

// 	return r
// }
