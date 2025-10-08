package main

import (
	"time"
)
import "github.com/google/uuid"

type Event struct {
	ID          *uuid.UUID `json:"id"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"` // Optional field
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	CreatedAt   *time.Time `json:"created_at"`
}

func (e *Event) ValidateCreation() *string {
	if e.Title == nil || *e.Title == "" || len(*e.Title) > 100 {
		errorMessage := `{"error": "Title is required and must be less than or equal to 100 characters"}`
		return &errorMessage
	}
	if e.StartTime == nil || e.EndTime == nil {
		errorMessage := `{"error": "Start time and End time are required"}`
		return &errorMessage
	}
	if e.StartTime.After(*e.EndTime) {
		errorMessage := `{"error": "Start time must be before end time"}`
		return &errorMessage
	}
	return nil
}
