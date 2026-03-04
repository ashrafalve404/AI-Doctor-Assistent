package handlers

import (
	"ai-doctor-bd/db"
	"ai-doctor-bd/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateProfile(w http.ResponseWriter, r *http.Request) {
	var p models.Profile
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec(
		"INSERT INTO profiles (name, age, gender, blood_group) VALUES (?, ?, ?, ?)",
		p.Name, p.Age, p.Gender, p.BloodGroup,
	)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func GetProfiles(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name, age, gender, blood_group, created_at FROM profiles ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	profiles := []models.Profile{}
	for rows.Next() {
		var p models.Profile
		rows.Scan(&p.ID, &p.Name, &p.Age, &p.Gender, &p.BloodGroup, &p.CreatedAt)
		profiles = append(profiles, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profiles)
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	db.DB.Exec("DELETE FROM messages WHERE consultation_id IN (SELECT id FROM consultations WHERE profile_id = ?)", id)
	db.DB.Exec("DELETE FROM consultations WHERE profile_id = ?", id)
	db.DB.Exec("DELETE FROM profiles WHERE id = ?", id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
