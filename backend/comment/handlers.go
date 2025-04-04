package handlers

import (
	"backend/db"
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {

	var comment models.Comments
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "invalid request ", http.StatusBadRequest)
		return
	}

	db.DB.Create(&comment)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	db.DB.Delete(&models.Comments{}, id)
	w.WriteHeader(http.StatusNoContent)

}

func GetComments(w http.ResponseWriter, r *http.Request) {

	var comments []models.Comments

	db.DB.Find(&comments)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&comments)
}

func GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var comment models.Comments
	result := db.DB.First(&comment, id)
	if result.Error != nil {
		http.Error(w, "No such comment ", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

func EditComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var comment models.Comments
	result := db.DB.First(&comment, id)
	if result.Error != nil {
		http.Error(w, "comment not found", http.StatusNotFound)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db.DB.Save(&comment)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)

}
