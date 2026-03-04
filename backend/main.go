package main

import (
	"ai-doctor-bd/db"
	"ai-doctor-bd/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load .env
	godotenv.Load()

	// Init DB
	db.Init()

	// Router
	r := mux.NewRouter()

	// Profile routes
	r.HandleFunc("/api/profiles", handlers.GetProfiles).Methods("GET")
	r.HandleFunc("/api/profiles", handlers.CreateProfile).Methods("POST")
	r.HandleFunc("/api/profiles/{id}", handlers.DeleteProfile).Methods("DELETE")

	// Consultation routes
	r.HandleFunc("/api/analyze", handlers.Analyze).Methods("POST")
	r.HandleFunc("/api/chat", handlers.ChatFollowUp).Methods("POST")
	r.HandleFunc("/api/history/{profileId}", handlers.GetHistory).Methods("GET")
	r.HandleFunc("/api/messages/{consultationId}", handlers.GetMessages).Methods("GET")

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	handler := c.Handler(r)

	log.Println("🚀 AI Doctor BD backend running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
