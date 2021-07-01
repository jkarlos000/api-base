package entity

// Session represents a admin record.
type Session struct {
	ID             string `json:"id"`
	Owner         string `json:"owner"`
	Diagram			string `json:"diagram"`
	Slug			string `json:"slug"`
	Password		string `json:"password"`
	IsActive       bool   `json:"is_active"`
}

// TableName represents the table name in the database
func (m Session) TableName() string {
	return "sessions"
}

