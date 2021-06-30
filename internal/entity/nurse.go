package entity

// Nurse represents a admin record.
type Nurse struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	HaveFullAccess bool   `json:"have_full_access"`
	IsActive       bool   `json:"is_active"`
	Username       string `json:"username" db:"username"`
	FirstName      string `json:"first_name" db:"first_name"`
	LastName       string `json:"last_name" db:"last_name"`
}

// TableName represents the table name in the database
func (m Nurse) TableName() string {
	return "nurses"
}

