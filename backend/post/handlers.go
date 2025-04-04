package post

import (
	"backend/db"
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./post-frontend/index.html")
}

func CreatePost(w http.ResponseWriter, r *http.Request) {

	var post models.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "invalid request ", http.StatusBadRequest)
		return
	}

	db.DB.Create(&post)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post

	db.DB.Find(&posts)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetPostByID(w http.ResponseWriter, r *http.Request) {
	/*
		mux (short for multiplexer) is a powerful HTTP request router and dispatcher for Go. It is commonly used in web applications to:
		Handle dynamic routes (e.g., /users/{id})
		Extract URL parameters (e.g., id from /users/5)
		Support middleware and request filtering
	*/
	/*
		When a request is routed through mux, mux.Vars(r) extracts the URL parameters and returns them as a map[string]string. This allows us to retrieve values from dynamic URL paths
	*/

	// Extract route parameters
	params := mux.Vars(r) //params --> map[string][string]

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var post models.Post
	result := db.DB.First(&post, id)

	if result.Error != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {

	//get post id
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	db.DB.Delete(&models.Post{}, id)
	w.WriteHeader(http.StatusNoContent)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "invalid post ID ", http.StatusBadRequest)
		return
	}

	var post models.Post
	result := db.DB.First(&post, id)
	if result.Error != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db.DB.Save(&post)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)

}
