package entity

import "time"

// Session represents a admin record.
type Session struct {
	ID             string `json:"id"`
	Owner         string `json:"owner"`
	Diagram			string `json:"diagram"`
	Slug			string `json:"slug"`
	Password		string `json:"password"`
	IsActive       bool   `json:"is_active"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time   `json:"updated_at" db:"updated_at"`
	DeletedAt	*time.Time `json:"deleted_at" db:"deleted_at"`
}

// TableName represents the table name in the database
func (m Session) TableName() string {
	return "sessions"
}

