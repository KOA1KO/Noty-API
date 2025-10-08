package models

import "time"

const (
	ASC  string = "ASC"
	DESC string = "DESC"
)

// Note представляет элемент выдачи.
// swagger:model
type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteMeta struct {
	Limit  int
	Offset int
	Sort   string
}
