package handlers

import (
	"ai-doctor-bd/ai"
	"ai-doctor-bd/db"
	"ai-doctor-bd/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Analyze(w http.ResponseWriter, r *http.Request) {
	var req models.SymptomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get profile info
	var p models.Profile
	err := db.DB.QueryRow("SELECT id, name, age, gender, blood_group FROM profiles WHERE id = ?", req.ProfileID).
		Scan(&p.ID, &p.Name, &p.Age, &p.Gender, &p.BloodGroup)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	profileInfo := fmt.Sprintf("নাম: %s, বয়স: %d, লিঙ্গ: %s, রক্তের গ্রুপ: %s", p.Name, p.Age, p.Gender, p.BloodGroup)

	// Call Qwen3 via Bytez
	response, urgency, err := ai.AnalyzeSymptoms(req.Symptoms, profileInfo)
	if err != nil {
		http.Error(w, "AI error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save consultation
	result, err := db.DB.Exec(
		"INSERT INTO consultations (profile_id, symptoms, ai_response, urgency) VALUES (?, ?, ?, ?)",
		req.ProfileID, req.Symptoms, response, urgency,
	)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	consultationID, _ := result.LastInsertId()

	// Save initial messages
	db.DB.Exec("INSERT INTO messages (consultation_id, role, content) VALUES (?, ?, ?)", consultationID, "user", req.Symptoms)
	db.DB.Exec("INSERT INTO messages (consultation_id, role, content) VALUES (?, ?, ?)", consultationID, "assistant", response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"consultation_id": consultationID,
		"response":        response,
		"urgency":         urgency,
	})
}

func ChatFollowUp(w http.ResponseWriter, r *http.Request) {
	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get consultation + profile
	var c models.Consultation
	var profileID int
	err := db.DB.QueryRow("SELECT id, profile_id, symptoms FROM consultations WHERE id = ?", req.ConsultationID).
		Scan(&c.ID, &profileID, &c.Symptoms)
	if err != nil {
		http.Error(w, "Consultation not found", http.StatusNotFound)
		return
	}

	var p models.Profile
	db.DB.QueryRow("SELECT name, age, gender, blood_group FROM profiles WHERE id = ?", profileID).
		Scan(&p.Name, &p.Age, &p.Gender, &p.BloodGroup)

	profileInfo := fmt.Sprintf("নাম: %s, বয়স: %d, লিঙ্গ: %s, রক্তের গ্রুপ: %s", p.Name, p.Age, p.Gender, p.BloodGroup)

	// Get chat history
	rows, _ := db.DB.Query("SELECT role, content FROM messages WHERE consultation_id = ? ORDER BY created_at ASC", req.ConsultationID)
	defer rows.Close()

	var history []ai.Message
	for rows.Next() {
		var msg ai.Message
		rows.Scan(&msg.Role, &msg.Content)
		history = append(history, msg)
	}

	// Call AI
	response, err := ai.Chat(req.Message, history, profileInfo, c.Symptoms)
	if err != nil {
		http.Error(w, "AI error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save messages
	db.DB.Exec("INSERT INTO messages (consultation_id, role, content) VALUES (?, ?, ?)", req.ConsultationID, "user", req.Message)
	db.DB.Exec("INSERT INTO messages (consultation_id, role, content) VALUES (?, ?, ?)", req.ConsultationID, "assistant", response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

func GetHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID := vars["profileId"]

	rows, err := db.DB.Query(
		"SELECT id, symptoms, ai_response, urgency, created_at FROM consultations WHERE profile_id = ? ORDER BY created_at DESC",
		profileID,
	)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	consultations := []models.Consultation{}
	for rows.Next() {
		var c models.Consultation
		rows.Scan(&c.ID, &c.Symptoms, &c.AIResponse, &c.Urgency, &c.CreatedAt)
		consultations = append(consultations, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(consultations)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	consultationID := vars["consultationId"]

	rows, err := db.DB.Query(
		"SELECT id, role, content, created_at FROM messages WHERE consultation_id = ? ORDER BY created_at ASC",
		consultationID,
	)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	messages := []models.Message{}
	for rows.Next() {
		var m models.Message
		rows.Scan(&m.ID, &m.Role, &m.Content, &m.CreatedAt)
		messages = append(messages, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
