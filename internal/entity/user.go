package entity

import (
	"fmt"
	"time"
)

// User represents a user record.
type User struct {
	ID          string       `json:"id" db:"id"`
	Username    string       `json:"username" db:"username"`
	Password    string       `json:"password" db:"password"`
	Email		string 		 `json:"email" db:"email"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time   `json:"updated_at" db:"updated_at"`
	Roles       []string     `json:"roles" db:"-"`
	Permissions []Permission `json:"permissions" db:"-"`
	IsActive    bool         `json:"is_active" db:"is_active"`
	FirstName   string       `json:"first_name" db:"first_name"`
	LastName    string       `json:"last_name" db:"last_name"`
	Account     interface{}  `json:"account" db:"-"`
	RoleID      string       `json:"role_id" db:"-"`
}

// TableName represents the table name
func (u User) TableName() string {
	return "users"
}

// GetID returns the user ID.
func (u User) GetID() string {
	return u.ID
}

// GetUsername returns the user username.
func (u User) GetUsername() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// GetRoles returns the user ID.
func (u User) GetRoles() []string {
	return u.Roles
}

// HasRole compare roles
func (u User) HasRole(roles ...string) bool {
	for _, role := range roles {
		for _, uRole := range u.Roles {
			if role == uRole {
				return true
			}
		}
	}
	return false
}

// IsUserActive GetStatus returns the user status.
func (u User) IsUserActive() bool {
	return u.IsActive
}

func (u User) GetEmail() string {
	return u.Email
}
