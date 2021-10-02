package entity

// Role represents a product record.
type Role struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// TableName represents the table name
func (c Role) TableName() string {
	return "roles"
}
