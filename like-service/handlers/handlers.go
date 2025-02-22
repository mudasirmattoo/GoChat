package handlers

import (
	"encoding/json"
	"like-service/database"
	"like-service/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func AddLike(w http.ResponseWriter, r *http.Request) {
	var like models.Like
	err := json.NewDecoder(r.Body).Decode(&like)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	database.DB.Create(&like)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(like)
}

func GetLikesByPostID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID, err := strconv.Atoi(params["post_id"])
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	var likes []models.Like
	database.DB.Where("post_id = ?", postID).Find(&likes)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(likes)
}

func RemoveLike(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	database.DB.Delete(&models.Like{}, id)
	w.WriteHeader(http.StatusNoContent)
}
