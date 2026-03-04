package models

import "time"

type Profile struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Age        int       `json:"age"`
	Gender     string    `json:"gender"`
	BloodGroup string    `json:"blood_group"`
	CreatedAt  time.Time `json:"created_at"`
}

type Consultation struct {
	ID         int       `json:"id"`
	ProfileID  int       `json:"profile_id"`
	Symptoms   string    `json:"symptoms"`
	AIResponse string    `json:"ai_response"`
	Urgency    string    `json:"urgency"`
	CreatedAt  time.Time `json:"created_at"`
}

type Message struct {
	ID             int       `json:"id"`
	ConsultationID int       `json:"consultation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

type SymptomRequest struct {
	ProfileID int    `json:"profile_id"`
	Symptoms  string `json:"symptoms"`
}

type ChatRequest struct {
	ConsultationID int    `json:"consultation_id"`
	Message        string `json:"message"`
}
