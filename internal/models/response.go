package models

import "time"

type ResponseNoteID struct {
	NewNoteID int `json:"new_note_id"`
}

type NoteResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
