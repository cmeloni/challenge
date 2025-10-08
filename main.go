package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type RestApp struct {
	eventTable *EventTable
}

func (a *RestApp) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var req Event
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	if err := req.ValidateCreation(); err != nil {
		http.Error(w, *err, http.StatusBadRequest)
		return
	}
	a.eventTable.Create(&req, r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		return
	}
}

func (a *RestApp) ListEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := a.eventTable.List(r.Context())
	if err != nil {
		http.Error(w, *err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	errDec := json.NewEncoder(w).Encode(events)
	if errDec != nil {
		return
	}
}

func (a *RestApp) GetEventHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/events/"):]
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid event ID format"}`, http.StatusBadRequest)
		return
	}
	event, errDes := a.eventTable.GetEventById(eventID, r.Context())
	if errDes != nil {
		http.Error(w, *errDes, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)

}

func NewRestApp(db *sql.DB) *RestApp {
	return &RestApp{
		eventTable: &EventTable{DB: db},
	}
}

func main() {
	db := GetDB()
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatalf("Failed to ping database: %v\n", err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")

	restApp := NewRestApp(db)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /events", restApp.CreateEventHandler)
	mux.HandleFunc("GET /events", restApp.ListEventsHandler)
	mux.HandleFunc("GET /events/{id}", restApp.GetEventHandler)

	port := ":8080"
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
