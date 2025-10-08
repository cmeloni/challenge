package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func GetDB() *sql.DB {
	connStr := "postgresql://postgres:123456@localhost:15432/challenge?sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return db
}

type EventTable struct {
	DB *sql.DB
}

func (e *EventTable) Create(event *Event, rCtx context.Context) *string {
	if event.ID == nil {
		uu := uuid.New()
		event.ID = &uu
	}
	query := `INSERT INTO collection.events (id, title, description, start_time, end_time)
              VALUES ($1, $2, $3, $4, $5) RETURNING created_at`
	ctx, cancel := context.WithTimeout(rCtx, 5*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query,
		*event.ID, *event.Title, event.Description, *event.StartTime, *event.EndTime).Scan(&event.CreatedAt)
	if err != nil {
		log.Printf("Error inserting event: %v", err)
		errorMessage := `{"error": "Failed to create event"}`
		return &errorMessage
	}
	return nil
}

func (e *EventTable) List(rCtx context.Context) ([]Event, *string) {
	query := `SELECT id, title, description, start_time, end_time, created_at FROM collection.events ORDER BY start_time ASC`
	ctx, cancel := context.WithTimeout(rCtx, 5*time.Second)
	defer cancel()
	rows, err := e.DB.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Error getting events: %v", err)
		errorMessage := `{"error": "Failed to get events"}`
		return nil, &errorMessage
	}
	defer rows.Close()
	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.CreatedAt)
		if err != nil {
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating event rows: %v", err)
		errorMessage := `{"error": "Failed to interacting events"}`
		return nil, &errorMessage
	}
	return events, nil
}

func (e *EventTable) GetEventById(id uuid.UUID, rCtx context.Context) (*Event, *string) {
	query := `SELECT id, title, description, start_time, end_time, created_at FROM collection.events WHERE id = $1`

	ctx, cancel := context.WithTimeout(rCtx, 5*time.Second)
	defer cancel()

	var event Event
	var description sql.NullString
	err := e.DB.QueryRowContext(ctx, query, id).Scan(
		&event.ID, &event.Title, &description, &event.StartTime, &event.EndTime, &event.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			errDesc := `{"error": "Event not found"}`
			return nil, &errDesc
		}
		log.Printf("Error querying event by ID: %v", err)
		errDesc := `{"error": "Failed to retrieve event"}`
		return nil, &errDesc
	}
	if description.Valid {
		event.Description = &description.String
	}
	return &event, nil
}
